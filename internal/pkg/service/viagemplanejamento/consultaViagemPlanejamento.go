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
)

var cfgInitiated bool
var logger logging.Logger
var loggerConcorrencia logging.Logger

//InitConfig - é responsável por iniciar configuração da package
func InitConfig() {
	if !cfgInitiated {
		logger = logging.NewLogger("service.viagemplanejamento", cfg.Config.Logging.Level)
		loggerConcorrencia = logging.NewLogger("service.viagemplanejamento.CONCORRENCIA", cfg.Config.Logging.Level)
		cfgInitiated = true
	}
}

//Service -
type Service struct {
	planEscRep   *repository.PlanejamentoEscalaRepository
	vigExecRep   *repository.ViagemExecutadaRepository
	cacheCliente *cache.Cliente

	err                        error
	chTotal                    chan int
	total                      int
	filaTrabalho               chan dto.FilterDTO
	resultado                  chan *dto.ConsultaViagemPlanejamentoDTO
	captura                    chan error
	wg                         sync.WaitGroup
	initiated                  chan bool
	consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO
	confirm                    uint8
}

//NewViagemPlanejamentoService -
func NewViagemPlanejamentoService(planEscRep *repository.PlanejamentoEscalaRepository, vigExecRep *repository.ViagemExecutadaRepository, cacheCliente *cache.Cliente) *Service {
	vps := &Service{}
	vps.planEscRep = planEscRep
	vps.vigExecRep = vigExecRep
	vps.cacheCliente = cacheCliente

	vps.filaTrabalho = make(chan dto.FilterDTO, 50)
	vps.resultado = make(chan *dto.ConsultaViagemPlanejamentoDTO, 5)
	vps.captura = make(chan error, 5)

	vps.initiated = make(chan bool, cfg.Config.Service.ViagemPlanejamento.MaxConcurrentSubTask*2)
	vps.chTotal = make(chan int, 1)
	go func() {
		vps.confirm = 0
		cInitiated := 0
		for {
			select {
			case resultadoParceialConsulta := <-vps.resultado:
				// consultaViagemPlanejamento.ViagensExecutada = append(consultaViagemPlanejamento.ViagensExecutada, resultadoParceialConsulta.ViagensExecutada...)
				vps.consultaViagemPlanejamento.ViagensExecutadaPendentes = append(vps.consultaViagemPlanejamento.ViagensExecutadaPendentes, resultadoParceialConsulta.ViagensExecutadaPendentes...)
				vps.consultaViagemPlanejamento.Viagens = append(vps.consultaViagemPlanejamento.Viagens, resultadoParceialConsulta.Viagens...)
				vps.confirm++
				loggerConcorrencia.Debugf("Confirm Ok [%d/%d]", vps.confirm, vps.total)
				vps.wg.Done()
			case vps.err = <-vps.captura:
				vps.confirm++
				loggerConcorrencia.Debugf("Confirm Err [%d/%d]", vps.confirm, vps.total)
				vps.wg.Done()
			case <-vps.initiated:
				cInitiated++
				loggerConcorrencia.Tracef("Initiated [%d/%d] ", cInitiated, vps.total)
			case vps.total = <-vps.chTotal:
				cInitiated = 0
			}
		}
	}()

	var processarConsultas = func() {
		for f := range vps.filaTrabalho {
			vps.initiated <- true
			vps.ConsultarPorTrajeto(f, vps.resultado, vps.captura)
		}
	}

	for i := 0; i < cfg.Config.Service.ViagemPlanejamento.MaxConcurrentSubTask; i++ {
		go processarConsultas()
	}

	return vps
}

//Consultar -
func (vps *Service) Consultar(filtro dto.FilterDTO) (*dto.ConsultaViagemPlanejamentoDTO, error) {
	start := time.Now()

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
				ListaTrajetos: []dto.TrajetoDTO{
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

	loggerConcorrencia.Debugf("ViagemPlanejamento.MaxConcurrentSubTask [%d] ", cfg.Config.Service.ViagemPlanejamento.MaxConcurrentSubTask)

	total := len(filtrosConsulta)
	vps.chTotal <- total
	loggerConcorrencia.Debugf("Pending [%d]", total)

	vps.wg = sync.WaitGroup{}
	vps.wg.Add(total)

	vps.consultaViagemPlanejamento = new(dto.ConsultaViagemPlanejamentoDTO)
	vps.consultaViagemPlanejamento.ViagensExecutada = []*model.ViagemExecutada{}
	vps.consultaViagemPlanejamento.Totalizadores = &dto.TotalizadoresDTO{}
	vps.consultaViagemPlanejamento.Viagens = []*dto.ViagemDTO{}

	enqueued := 0
	for _, f := range filtrosConsulta {
		enqueued++
		loggerConcorrencia.Tracef("Enqueued [%d/%d] ", enqueued, total)
		vps.filaTrabalho <- f
	}

	vps.wg.Wait()

	dto.OrdenarViagemPorData(vps.consultaViagemPlanejamento.Viagens)

	processarAtrasadas(vps.consultaViagemPlanejamento)

	calcularTotalizadores(vps.consultaViagemPlanejamento)

	// for _, vg := range vps.consultaViagemPlanejamento.Viagens {
	// 	fmt.Printf("%v - %v\n", vg.DataAbertura, vg.IDViagemExecutada.Hex())
	// }

	duracao := time.Since(start)

	var informacoes = make(map[string]interface{})
	informacoes["duracao"] = fmt.Sprintf("%v", duracao)
	vps.consultaViagemPlanejamento.Informacoes = informacoes

	logger.Debugf("QTD Total de Viagens: %d\t em %v\n", len(vps.consultaViagemPlanejamento.Viagens), duracao)

	return vps.consultaViagemPlanejamento, vps.err
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
				if vg.Status == dto.StatusViagem.NaoRealizada && vg.IDViagemExecutada == "" {
					vg.Status = dto.StatusViagem.Atrasada
					populaDadosViagem(vgex, vg)
				} else {
					novaVG := &dto.ViagemDTO{}
					novaVG.Status = dto.StatusViagem.Extra
					populaDadosViagem(vgex, novaVG)
					vgRealizadasNaoPlanejadas = append(vgRealizadasNaoPlanejadas, novaVG)
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
	Extra                chan int32
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
	close(tot.Extra)
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
	lancar(acumulador, &tot.Extra, &t.Extra)
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

	if vg.Status == dto.StatusViagem.RealizadaPlanejada || vg.Status == dto.StatusViagem.Atrasada || vg.Status == dto.StatusViagem.Extra {
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

	if vg.Status == dto.StatusViagem.Extra {
		t.Extra <- 1
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
	var errPlanejamentos error
	var errViagensExecutadas error

	consultaViagemPlanejamentoDTO := &dto.ConsultaViagemPlanejamentoDTO{}
	viagensDTO := []*dto.ViagemDTO{}
	viagensExecutadaPendentes := []*model.ViagemExecutada{}

	retornoMapaHorarioViagem := make(chan map[int32]*dto.ViagemDTO)

	go func() {
		mapaHorarioViagemAUX := make(map[int32]*dto.ViagemDTO)
		var planejamentosEscala []*model.ProcPlanejamentoEscala
		planejamentosEscala, errPlanejamentos = vps.planEscRep.ListarPlanejamentosEscala(&filtro)
		if errPlanejamentos != nil {
			logger.Errorf("Erro ao ListarPlanejamentosEscala - %s\n", errPlanejamentos)
			retornoMapaHorarioViagem <- nil
			return
		}

		for _, ples := range planejamentosEscala {
			vg, _ := converterPlanejamentosEscala(ples, filtro)
			viagensDTO = append(viagensDTO, vg)
			mapaHorarioViagemAUX[vg.IDHorario] = vg
		}
		retornoMapaHorarioViagem <- mapaHorarioViagemAUX
	}()

	viagensExecutada, errViagensExecutadas := vps.vigExecRep.ListarViagensPor(filtro)
	if errViagensExecutadas != nil {
		logger.Errorf("Erro ao ListarViagensPor %+v - %s\n", filtro, errViagensExecutadas)
	}

	mapaHorarioViagem := <-retornoMapaHorarioViagem

	if errPlanejamentos != nil || errViagensExecutadas != nil {
		var err error
		if errPlanejamentos != nil {
			err = errPlanejamentos
		}
		if errViagensExecutadas != nil {
			if err != nil {
				err = fmt.Errorf("%v\n%v", err, errViagensExecutadas)
			} else {
				err = errViagensExecutadas

			}
		}
		captura <- err
		return
	}

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

func converterPlanejamentosEscala(ples *model.ProcPlanejamentoEscala, filtro dto.FilterDTO) (*dto.ViagemDTO, error) {
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
		Trajeto: dto.TrajetoDTO{
			ID:        filtro.ListaTrajetos[0].ID,
			Descricao: filtro.ListaTrajetos[0].Descricao,
			Sentido:   filtro.ListaTrajetos[0].Sentido,
		},
	}

	if ples.Partida.Before(filtro.Complemento.DataHora) {
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
	vg.Trajeto = dto.TrajetoDTO{
		ID:        vgex.Partida.TrajetoExecutado.IDObject,
		Descricao: vgex.Partida.TrajetoExecutado.Descricao,
		Sentido:   vgex.Partida.TrajetoExecutado.Sentido,
	}

}

func diferencaMinutos(d time.Duration) (duracao int64) {
	diff := d.Minutes()
	if diff >= 1 || diff <= -1 {
		duracao = int64(diff)
	}
	return
}
