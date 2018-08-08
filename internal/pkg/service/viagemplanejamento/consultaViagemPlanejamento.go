package viagemplanejamento

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/cache"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/util"
	"gopkg.in/mgo.v2/bson"
)

var logger logging.Logger
var loggerConcorrencia logging.Logger

//InitConfig - é responsável por iniciar configuração da package
func InitConfig() {
	logger = logging.NewLogger("service.viagemplanejamento", cfg.Config.Logging.Level)
	loggerConcorrencia = logging.NewLogger("service.viagemplanejamento.CONCORRENCIA", cfg.Config.Logging.Level)
}

//Service -
type Service struct {
	planEscRep   *repository.PlanejamentoEscalaRepository
	vigExecRep   *repository.ViagemExecutadaRepository
	cacheCliente *cache.Cliente
}

//NewViagemPlanejamentoService -
func NewViagemPlanejamentoService(planEscRep *repository.PlanejamentoEscalaRepository, vigExecRep *repository.ViagemExecutadaRepository, cacheCliente *cache.Cliente) *Service {
	vps := &Service{}
	vps.planEscRep = planEscRep
	vps.vigExecRep = vigExecRep
	vps.cacheCliente = cacheCliente
	return vps
}

//Consultar -
func (vps *Service) Consultar(filtro dto.FilterDTO) (*dto.ConsultaViagemPlanejamentoDTO, error) {
	start := time.Now()
	var err error
	cliente := vps.cacheCliente.Cache[filtro.IDCliente]
	filtro.Complemento = dto.DadosComplementares{
		Cliente:  cliente,
		DataHora: time.Now(),
	}

	periodo := util.Periodo{Inicio: filtro.GetDataInicio(), Fim: filtro.GetDataFim()}
	periodos := util.SplitDiasPeriodo(periodo)
	trajetos := filtro.ListaTrajetos
	var filtrosConsulta []dto.FilterDTO
	for _, p := range periodos {
		for _, t := range trajetos {
			f := dto.FilterDTO{
				ListaTrajetos: []bson.ObjectId{
					t,
				},
				IDCliente:   filtro.IDCliente,
				IDVeiculo:   filtro.IDVeiculo,
				Ordenacao:   filtro.Ordenacao,
				DataInicio:  util.FormatarAMDHMS(p.Inicio),
				DataFim:     util.FormatarAMDHMS(p.Fim),
				TipoDia:     model.TiposDia.FromDate(p.Inicio, []string{"O", "F"}),
				Complemento: filtro.Complemento,
			}
			filtrosConsulta = append(filtrosConsulta, f)
		}
	}

	loggerConcorrencia.Debugf("ViagemPlanejamento.MaxConcurrent [%d] ", cfg.Config.Service.ViagemPlanejamento.MaxConcurrent)

	total := len(filtrosConsulta)
	loggerConcorrencia.Debugf("Pending [%d]", total)

	//TODO - Tornar channels atributos de instância para diminuir quantidade de objetos criados, com isso
	//		Channel concluido não será mais necessário
	// Criar pool de *viagemplanejamento.Service para limitar quantidade de consultas simultâneas
	filaTrabalho := make(chan dto.FilterDTO, 50)
	resultado := make(chan *dto.ConsultaViagemPlanejamentoDTO, 5)
	captura := make(chan error, 5)
	concluido := make(chan bool, 2)
	defer close(filaTrabalho)
	defer close(resultado)
	defer close(captura)
	defer close(concluido)

	var wg sync.WaitGroup
	wg.Add(total)

	initiated := 0
	var processarConsultas = func() {
		for f := range filaTrabalho {
			initiated++
			loggerConcorrencia.Tracef("Initiated [%d/%d] ", initiated, total)
			vps.ConsultarPorTrajeto(f, resultado, captura)
		}
	}

	for i := 0; i < cfg.Config.Service.ViagemPlanejamento.MaxConcurrent; i++ {
		go processarConsultas()
	}

	consultaViagemPlanejamento := new(dto.ConsultaViagemPlanejamentoDTO)
	consultaViagemPlanejamento.ViagensExecutada = []*model.ViagemExecutada{}
	consultaViagemPlanejamento.Totalizadores = &dto.TotalizadoresDTO{}
	consultaViagemPlanejamento.Viagens = []*dto.ViagemDTO{}

	go func() {
		confirm := 0
		for {
			var b bool
			select {
			case resultadoParceialConsulta, b := <-resultado:
				if b {
					// consultaViagemPlanejamento.ViagensExecutada = append(consultaViagemPlanejamento.ViagensExecutada, resultadoParceialConsulta.ViagensExecutada...)
					consultaViagemPlanejamento.ViagensExecutadaPendentes = append(consultaViagemPlanejamento.ViagensExecutadaPendentes, resultadoParceialConsulta.ViagensExecutadaPendentes...)
					consultaViagemPlanejamento.Viagens = append(consultaViagemPlanejamento.Viagens, resultadoParceialConsulta.Viagens...)
					confirm++
					loggerConcorrencia.Debugf("Confirm Ok [%d/%d]", confirm, total)
					wg.Done()
				}
			case err, b = <-captura:
				if b {
					confirm++
					loggerConcorrencia.Debugf("Confirm Err [%d/%d]", confirm, total)
					wg.Done()
				}
			case <-concluido:
				return
			}
		}
	}()

	enqueued := 0
	for _, f := range filtrosConsulta {
		enqueued++
		loggerConcorrencia.Tracef("Enqueued [%d/%d] ", enqueued, total)
		filaTrabalho <- f
	}

	wg.Wait()

	dto.OrdenarViagemPorData(consultaViagemPlanejamento.Viagens)

	processarAtrasadas(consultaViagemPlanejamento)

	calcularTotalizadores(consultaViagemPlanejamento)

	duracao := time.Since(start)

	var informacoes = make(map[string]interface{})
	informacoes["duracao"] = fmt.Sprintf("%v", duracao)
	consultaViagemPlanejamento.Informacoes = informacoes

	logger.Debugf("QTD Total de Viagens: %d\t em %v\n", len(consultaViagemPlanejamento.Viagens), duracao)

	concluido <- true
	return consultaViagemPlanejamento, err
}

func processarAtrasadas(consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO) {
	model.OrdenarViagemExecutadaPorData(consultaViagemPlanejamento.ViagensExecutadaPendentes)

	vgRealizadasNaoPlanejadas := []*dto.ViagemDTO{}

	pos := 0
	size := len(consultaViagemPlanejamento.Viagens)

	for _, vgex := range consultaViagemPlanejamento.ViagensExecutadaPendentes {
		for ; pos < size; pos++ {
			j := pos + 1

			vg := consultaViagemPlanejamento.Viagens[pos]
			if j < size && vgex.Executada.DataInicio.After(consultaViagemPlanejamento.Viagens[j].PartidaOrdenacao) {
				continue
			}

			if vgex.Executada.DataInicio.After(vg.PartidaOrdenacao) {
				if vg.Status == dto.StatusViagem.NaoRealizada {
					if vg.IDViagemExecutada == "" {
						vg.Status = dto.StatusViagem.Atrasada
						populaDadosViagem(vgex, vg)
					} else {
						novaVG := &dto.ViagemDTO{}
						novaVG.Status = dto.StatusViagem.RealizadaNaoPlanejada
						populaDadosViagem(vgex, novaVG)
						vgRealizadasNaoPlanejadas = append(vgRealizadasNaoPlanejadas, novaVG)
					}
				}
				break
			}
		}
	}

	consultaViagemPlanejamento.Viagens = append(consultaViagemPlanejamento.Viagens, vgRealizadasNaoPlanejadas...)
	dto.OrdenarViagemPorData(consultaViagemPlanejamento.Viagens)
}

func calcularTotalizadores(consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO) {

	var wgTot *sync.WaitGroup
	wgTot = &sync.WaitGroup{}

	tot, tots := newTotalizacao(consultaViagemPlanejamento.Totalizadores, wgTot)

	chTotVG := make(chan *dto.ViagemDTO)
	totalizarParalelamente := func() {
		for vg := range chTotVG {
			totalizar(vg, tot, wgTot)
		}
	}

	for i := 0; i < 10; i++ {
		go totalizarParalelamente()
	}

	for _, vg := range consultaViagemPlanejamento.Viagens {
		wgTot.Add(tots)
		chTotVG <- vg
	}

	wgTot.Wait()

	consultaViagemPlanejamento.Totalizadores.NaoIniciadas = (consultaViagemPlanejamento.Totalizadores.Planejadas - consultaViagemPlanejamento.Totalizadores.PlanejadasAteMomento)

	indiceExecucao := (float64(consultaViagemPlanejamento.Totalizadores.Realizadas) / float64(consultaViagemPlanejamento.Totalizadores.PlanejadasAteMomento) * 100)
	consultaViagemPlanejamento.Totalizadores.IndiceExecucao = []int32{int32(indiceExecucao)}

	defer tot.close()
	close(chTotVG)

}

type totalizacao struct {
	Planejadas           chan int32
	Realizadas           chan int32
	RealizadasPlanejadas chan int32
	EmAndamento          chan int32
	Canceladas           chan int32
	NaoIniciadas         chan int32
	NaoRealizadas        chan int32
	Atrasada             chan int32
	PlanejadasAteMomento chan int32
}

func (tot *totalizacao) close() {
	close(tot.Planejadas)
	close(tot.Realizadas)
	close(tot.RealizadasPlanejadas)
	close(tot.EmAndamento)
	close(tot.Canceladas)
	close(tot.NaoIniciadas)
	close(tot.NaoRealizadas)
	close(tot.Atrasada)
	close(tot.PlanejadasAteMomento)
}

func newTotalizacao(t *dto.TotalizadoresDTO, wg *sync.WaitGroup) (tot *totalizacao, tots int) {
	chSize := 10
	tot = &totalizacao{}

	lancar := func(f func(ch *chan int32, p *int32), ch *chan int32, p *int32) {
		tots++
		*ch = make(chan int32, chSize)
		go f(ch, p)
	}

	acumulador := func(ch *chan int32, p *int32) {
		for v := range *ch {
			*p += v
			wg.Done()
		}
	}

	lancar(acumulador, &tot.Planejadas, &t.Planejadas)
	lancar(acumulador, &tot.Realizadas, &t.Realizadas)
	lancar(acumulador, &tot.RealizadasPlanejadas, &t.RealizadasPlanejadas)
	lancar(acumulador, &tot.EmAndamento, &t.EmAndamento)
	lancar(acumulador, &tot.Canceladas, &t.Canceladas)
	lancar(acumulador, &tot.NaoIniciadas, &t.NaoIniciadas)
	lancar(acumulador, &tot.NaoRealizadas, &t.NaoRealizadas)
	lancar(acumulador, &tot.Atrasada, &t.Atrasada)
	lancar(acumulador, &tot.PlanejadasAteMomento, &t.PlanejadasAteMomento)

	return
}

func totalizar(vg *dto.ViagemDTO, t *totalizacao, wg *sync.WaitGroup) {

	if vg.Planejada {
		t.Planejadas <- 1
	} else {
		wg.Done()
	}

	if vg.Status == dto.StatusViagem.RealizadaPlanejada || vg.Status == dto.StatusViagem.Atrasada {
		t.Realizadas <- 1
	} else {
		wg.Done()
	}

	if vg.Status == dto.StatusViagem.RealizadaPlanejada {
		t.RealizadasPlanejadas <- 1
	} else {
		wg.Done()
	}

	if vg.Status == dto.StatusViagem.NaoIniciada {
		t.NaoIniciadas <- 1
	} else {
		wg.Done()
	}

	if vg.Status == dto.StatusViagem.EmAndamento {
		t.EmAndamento <- 1
	} else {
		wg.Done()
	}

	if vg.Status == dto.StatusViagem.Cancelada {
		t.Canceladas <- 1
	} else {
		wg.Done()
	}

	if vg.Status == dto.StatusViagem.NaoRealizada {
		t.NaoRealizadas <- 1
	} else {
		wg.Done()
	}

	if vg.Status == dto.StatusViagem.Atrasada {
		t.Atrasada <- 1
	} else {
		wg.Done()
	}

	if vg.PlanejadaAteMomento {
		t.PlanejadasAteMomento <- 1
	} else {
		wg.Done()
	}
}

//ConsultarPorTrajeto -
func (vps *Service) ConsultarPorTrajeto(filtro dto.FilterDTO, resultado chan *dto.ConsultaViagemPlanejamentoDTO, captura chan error) {
	var err error

	consultaViagemPlanejamentoDTO := &dto.ConsultaViagemPlanejamentoDTO{}
	viagensDTO := []*dto.ViagemDTO{}
	viagensExecutadaPendentes := []*model.ViagemExecutada{}

	retornoMapaHorarioViagem := make(chan map[int32]*dto.ViagemDTO)

	go func() {
		mapaHorarioViagemAUX := make(map[int32]*dto.ViagemDTO)
		planejamentosEscala, err := vps.planEscRep.ListarPlanejamentosEscala(&filtro)
		if err != nil {
			logger.Errorf("Erro ao ListarPlanejamentosEscala - %s\n", err)
		}

		for _, ples := range planejamentosEscala {
			vg, _ := converterPlanejamentosEscala(ples, filtro.Complemento.DataHora)
			viagensDTO = append(viagensDTO, vg)
			mapaHorarioViagemAUX[vg.IDHorario] = vg
		}
		retornoMapaHorarioViagem <- mapaHorarioViagemAUX
	}()

	viagensExecutada, err := vps.vigExecRep.ListarViagensPor(filtro)
	if err != nil {
		logger.Errorf("Erro ao ListarViagensPor %+v - %s\n", filtro, err)
		captura <- err
		return
	}

	mapaHorarioViagem := <-retornoMapaHorarioViagem
	for _, vgex := range viagensExecutada {
		vgexNaoEncontrada, vg, _ := converterViagemExecutada(vgex, mapaHorarioViagem)
		if vg != nil {
			viagensDTO = append(viagensDTO, vg)

		} else if vgexNaoEncontrada != nil {

			viagensExecutadaPendentes = append(viagensExecutadaPendentes, vgexNaoEncontrada)
		}

	}

	// consultaViagemPlanejamentoDTO.ViagensExecutada = viagensExecutada
	consultaViagemPlanejamentoDTO.Viagens = viagensDTO
	consultaViagemPlanejamentoDTO.ViagensExecutadaPendentes = viagensExecutadaPendentes

	logger.Tracef("%+v\n", consultaViagemPlanejamentoDTO)

	resultado <- consultaViagemPlanejamentoDTO
}

func converterPlanejamentosEscala(ples *model.ProcPlanejamentoEscala, dataHora time.Time) (*dto.ViagemDTO, error) {
	var vg *dto.ViagemDTO
	var err error

	logger.Tracef("%+v\n", ples)
	// logger.Tracef("%#v\n", ples)

	vg = &dto.ViagemDTO{
		IDTabela:           ples.IDTabela,
		IDHorario:          ples.IDHorario,
		IDEmpresaPlanejada: ples.IDEmpresaPlan,
		NmTabela:           ples.NmTabela,
		PartidaOrdenacao:   ples.Partida,
		PartidaPlanTime:    ples.Partida,
		PartidaPlan:        util.FormatarHMS(ples.Partida),
		ChegadaPlanTime:    ples.Chegada,
		ChegadaPlan:        util.FormatarHMS(ples.Chegada),
		VeiculoPlan:        strconv.Itoa(int(ples.CodVeiculoPlan)),
		Planejada:          true,
	}

	if ples.Partida.Before(dataHora) {
		vg.Status = dto.StatusViagem.NaoRealizada
		vg.PlanejadaAteMomento = true
	} else {
		vg.Status = dto.StatusViagem.NaoIniciada
	}

	return vg, err
}

func converterViagemExecutada(vgex *model.ViagemExecutada, mapaHorarioViagem map[int32]*dto.ViagemDTO) (naoEncontrada *model.ViagemExecutada, processada *dto.ViagemDTO, err error) {

	logger.Tracef("%+v\n", vgex)
	// logger.Tracef("%#v\n", vgex)
	var vg *dto.ViagemDTO

	if vgex.Alocacao.IDHorario != "" {
		idHorario, err := strconv.Atoi(vgex.Alocacao.IDHorario)
		if err == nil && idHorario > 0 {
			vg = mapaHorarioViagem[int32(idHorario)]
		}
		err = nil
	}

	if vgex.SituacaoAtual == model.ViagemEstado.ViagemCancelada {
		vg = &dto.ViagemDTO{}
		vg.Status = dto.StatusViagem.Cancelada
		processada = vg
	} else if vg != nil && vg.IDViagemExecutada == "" {
		vg.Status = dto.StatusViagem.RealizadaPlanejada
	} else {
		vg = nil
		naoEncontrada = &(*vgex)
		return
	}

	populaDadosViagem(vgex, vg)

	return
}

func populaDadosViagem(vgex *model.ViagemExecutada, vg *dto.ViagemDTO) {

	var vgexDataFim time.Time
	if vgex.Executada.DataFim == nil || vgex.SituacaoAtual == model.ViagemEstado.NovaViagem || vgex.SituacaoAtual == model.ViagemEstado.EmPreparacao || vgex.SituacaoAtual == model.ViagemEstado.ViagemAberta || vgex.SituacaoAtual == model.ViagemEstado.DeslocamentoEmCerca {
		vg.EmExecucao = true
		vg.Status = dto.StatusViagem.EmAndamento
		vgexDataFim = time.Time{}
	} else {
		vgexDataFim = *vgex.Executada.DataFim
		vg.ChegadaReal = util.FormatarHMS(vgexDataFim)
		vg.DataFechamento = util.FormatarAMDHMS(vgexDataFim)
		duracao, duracaoFormatada := util.DuracaoEFormatacao(vgex.Executada.DataInicio, vgexDataFim)
		vg.Duracao = duracaoFormatada
		vg.DuracaoSeg = diferencaMinutos(duracao)
	}

	// if vgex.SituacaoAtual == model.ViagemEstado.NovaViagem || vgex.SituacaoAtual == model.ViagemEstado.EmPreparacao || vgex.SituacaoAtual == model.ViagemEstado.ViagemAberta || vgex.SituacaoAtual == model.ViagemEstado.DeslocamentoEmCerca {
	// 	vg.EmExecucao = true
	// 	vg.Status = dto.StatusViagem.EmAndamento
	// 	vgexDataFim = time.Time{}
	// }
	// if vgex.Executada.DataFim == nil {
	// 	vgexDataFim = *vgex.Executada.DataFim
	// 	vg.ChegadaReal = util.FormatarHMS(vgexDataFim)
	// 	vg.DataFechamento = util.FormatarAMDHMS(vgexDataFim)
	// 	duracao, duracaoFormatada := util.DuracaoEFormatacao(vgex.Executada.DataInicio, vgexDataFim)
	// 	vg.Duracao = duracaoFormatada
	// 	vg.DuracaoSeg = diferencaMinutos(duracao)
	// }

	vg.IDViagemExecutada = vgex.ID
	vg.Ipk = vgex.Ipk
	if len(vgex.PorcentagemConclusao) > 0 {
		fmtPorcentagemConclusao := vgex.PorcentagemConclusao
		fmtPorcentagemConclusao = strings.Replace(fmtPorcentagemConclusao, ",", ".", -1)
		if len(strings.Split(fmtPorcentagemConclusao, ".")) == 1 {
			fmtPorcentagemConclusao += ".00"
		}
		vg.PercentualConclusao = fmtPorcentagemConclusao
	}

	if vgex.Executada.Veiculo.ID > 0 {
		vg.IDVeiculo = strconv.Itoa(int(vgex.Executada.Veiculo.ID))
	}
	vg.VeiculoReal = vgex.Executada.Veiculo.Prefixo

	vg.PartidaReal = util.FormatarHMS(vgex.Executada.DataInicio)

	vg.Data = util.FormatarAMDHMS(vgex.Executada.DataInicio)
	vg.DataAbertura = vg.Data

	if vg.IDHorario > 0 { //Se planejamento encontrado
		diffPartida, diffPartidaFormatada := util.DuracaoEFormatacaoMinutos(vg.PartidaPlanTime, vgex.Executada.DataInicio)
		vg.DiffPartidaStr = diffPartidaFormatada
		vg.DiffPartida = diferencaMinutos(diffPartida)

		if vgex.Executada.DataFim != nil {
			diffChegada, diffChegadaFormatada := util.DuracaoEFormatacaoMinutos(vg.ChegadaPlanTime, vgexDataFim)
			vg.DiffChegadaStr = diffChegadaFormatada
			vg.DiffChegada = diferencaMinutos(diffChegada)
		}
	} else { //Se planejamento não encontrado
		vg.PartidaOrdenacao = vgex.Executada.DataInicio
	}

	vg.VelocidadeMedia = vgex.VelocidadeMedia

	vg.QtdePassageiros = vgex.QntPassageiros

}

func diferencaMinutos(d time.Duration) (duracao int64) {
	diff := d.Minutes()
	if diff >= 1 || diff <= -1 {
		duracao = int64(diff)
	}
	return
}
