package viagemplanejamento

import (
	"fmt"
	"strconv"
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
	periodo := util.Periodo{Inicio: filtro.GetDataInicio(), Fim: filtro.GetDataFim()}
	periodos := util.SplitDiasPeriodo(periodo)
	trajetos := filtro.ListaTrajetos
	var filtrosConsulta []dto.FilterDTO
	cliente := vps.cacheCliente.Cache[filtro.IDCliente]
	for _, p := range periodos {
		for _, t := range trajetos {
			f := dto.FilterDTO{
				ListaTrajetos: []bson.ObjectId{
					t,
				},
				IDCliente:  filtro.IDCliente,
				IDVeiculo:  filtro.IDVeiculo,
				Ordenacao:  filtro.Ordenacao,
				DataInicio: util.FormatarAMDHMS(p.Inicio),
				DataFim:    util.FormatarAMDHMS(p.Fim),
				TipoDia:    model.TiposDia.FromDate(p.Inicio, []string{"O", "F"}),
				Complemento: dto.DadosComplementares{
					Cliente: cliente,
				},
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

	enqueued := 0
	for _, f := range filtrosConsulta {
		enqueued++
		loggerConcorrencia.Tracef("Enqueued [%d/%d] ", enqueued, total)
		filaTrabalho <- f
	}

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
					wg.Done()
					// consultaViagemPlanejamento.ViagensExecutada = append(consultaViagemPlanejamento.ViagensExecutada, resultadoParceialConsulta.ViagensExecutada...)
					consultaViagemPlanejamento.Viagens = append(consultaViagemPlanejamento.Viagens, resultadoParceialConsulta.Viagens...)
					confirm++
					loggerConcorrencia.Debugf("Confirm Ok [%d/%d]", confirm, total)

				}
			case err, b = <-captura:
				if b {
					wg.Done()
					confirm++
					loggerConcorrencia.Debugf("Confirm Err [%d/%d]", confirm, total)
				}
			case <-concluido:
				return
			}
		}
	}()
	wg.Wait()

	dto.OrdenarViagemExecutadaPorData(consultaViagemPlanejamento.Viagens)

	var totPlanejadas int32
	var totPlanejadasAteMomento int32
	var totRealizadas int32
	var totRealizadasPlanejadas int32
	var totEmAndamento int32
	var totCanceladas int32
	var totPassageiros int32
	var totNaoIniciadas int32
	var totNaoRealizadas int32
	var totReforco int32
	var totAtrasada int32

	totPlanejadas = 5
	totPlanejadasAteMomento = 5
	totRealizadas = 5
	totRealizadasPlanejadas = 5
	totEmAndamento = 5
	totCanceladas = 5
	totPassageiros = 5
	totNaoIniciadas = 5
	totNaoRealizadas = 5
	totReforco = 5
	totAtrasada = 5

	consultaViagemPlanejamento.Totalizadores.Planejadas = totPlanejadas
	consultaViagemPlanejamento.Totalizadores.PlanejadasAteMomento = totPlanejadasAteMomento
	consultaViagemPlanejamento.Totalizadores.Realizadas = totRealizadas
	consultaViagemPlanejamento.Totalizadores.RealizadasPlanejadas = totRealizadasPlanejadas
	consultaViagemPlanejamento.Totalizadores.EmAndamento = totEmAndamento
	consultaViagemPlanejamento.Totalizadores.Canceladas = totCanceladas
	consultaViagemPlanejamento.Totalizadores.Passageiros = totPassageiros
	consultaViagemPlanejamento.Totalizadores.NaoIniciadas = totNaoIniciadas
	consultaViagemPlanejamento.Totalizadores.NaoRealizadas = totNaoRealizadas
	consultaViagemPlanejamento.Totalizadores.Reforco = totReforco
	consultaViagemPlanejamento.Totalizadores.Atrasada = totAtrasada
	consultaViagemPlanejamento.Totalizadores.IndiceExecucao = []int32{int32(len(consultaViagemPlanejamento.Viagens))}

	duracao := time.Since(start)

	var informacoes = make(map[string]interface{})
	informacoes["duracao"] = fmt.Sprintf("%v", duracao)
	consultaViagemPlanejamento.Informacoes = informacoes

	logger.Debugf("QTD Total de Viagens: %d\t em %v\n", len(consultaViagemPlanejamento.Viagens), duracao)

	concluido <- true
	return consultaViagemPlanejamento, err
}

//ConsultarPorTrajeto -
func (vps *Service) ConsultarPorTrajeto(filtro dto.FilterDTO, resultado chan *dto.ConsultaViagemPlanejamentoDTO, captura chan error) {
	var err error

	consultaViagemPlanejamentoDTO := &dto.ConsultaViagemPlanejamentoDTO{}
	viagensDTO := []*dto.ViagemDTO{}

	retornoMapaHorarioViagem := make(chan map[int32]*dto.ViagemDTO)

	go func() {
		mapaHorarioViagemAUX := make(map[int32]*dto.ViagemDTO)
		planejamentosEscala, err := vps.planEscRep.ListarPlanejamentosEscala(&filtro)
		if err != nil {
			logger.Errorf("Erro ao ListarPlanejamentosEscala - %s\n", err)
		}

		for _, ples := range planejamentosEscala {
			vg, _ := converterPlanejamentosEscala(ples)
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

		}

	}

	// consultaViagemPlanejamentoDTO.ViagensExecutada = viagensExecutada
	consultaViagemPlanejamentoDTO.Viagens = viagensDTO

	logger.Tracef("%+v\n", consultaViagemPlanejamentoDTO)

	resultado <- consultaViagemPlanejamentoDTO
}

func converterPlanejamentosEscala(ples *model.ProcPlanejamentoEscala) (*dto.ViagemDTO, error) {
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
		Status:             dto.StatusViagem.NaoRealizada,
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

	vg.IDViagemExecutada = vgex.ID
	vg.Ipk = vgex.Ipk
	vg.PercentualConclusao = vgex.PorcentagemConclusao

	if vgex.Executada.Veiculo.ID > 0 {
		vg.IDVeiculo = strconv.Itoa(int(vgex.Executada.Veiculo.ID))
	}
	vg.VeiculoReal = vgex.Executada.Veiculo.Prefixo

	vg.PartidaReal = util.FormatarHMS(vgex.Executada.DataInicio)
	vg.ChegadaReal = util.FormatarHMS(vgex.Executada.DataFim)

	vg.Data = util.FormatarAMDHMS(vgex.Executada.DataInicio)
	vg.DataAbertura = vg.Data
	vg.DataFechamento = util.FormatarAMDHMS(vgex.Executada.DataFim)

	duracao, duracaoFormatada := util.DuracaoEFormatacao(vgex.Executada.DataInicio, vgex.Executada.DataFim)
	vg.Duracao = duracaoFormatada
	vg.DuracaoSeg = diferencaMinutos(duracao)

	if vg.IDHorario > 0 { //Se planejamento encontrado
		diffPartida, diffPartidaFormatada := util.DuracaoEFormatacaoMinutos(vg.PartidaPlanTime, vgex.Executada.DataInicio)
		vg.DiffPartidaStr = diffPartidaFormatada
		vg.DiffPartida = diferencaMinutos(diffPartida)

		diffChegada, diffChegadaFormatada := util.DuracaoEFormatacaoMinutos(vg.ChegadaPlanTime, vgex.Executada.DataFim)
		vg.DiffChegadaStr = diffChegadaFormatada
		vg.DiffChegada = diferencaMinutos(diffChegada)
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
