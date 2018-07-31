package viagemplanejamento

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
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
	planEscRep *repository.PlanejamentoEscalaRepository
	vigExecRep *repository.ViagemExecutadaRepository
}

//NewViagemPlanejamentoService -
func NewViagemPlanejamentoService(planEscRep *repository.PlanejamentoEscalaRepository, vigExecRep *repository.ViagemExecutadaRepository) *Service {
	vps := &Service{}
	vps.planEscRep = planEscRep
	vps.vigExecRep = vigExecRep
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
			}
			filtrosConsulta = append(filtrosConsulta, f)
		}
	}

	loggerConcorrencia.Debugf("ViagemPlanejamento.MaxConcurrent [%d] ", cfg.Config.Service.ViagemPlanejamento.MaxConcurrent)

	total := len(filtrosConsulta)
	loggerConcorrencia.Debugf("Pending [%d]", total)

	filaTrabalho := make(chan dto.FilterDTO, 50)
	resultado := make(chan *dto.ConsultaViagemPlanejamentoDTO, 5)
	captura := make(chan error, 5)

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
	consultaViagemPlanejamento.Viagens = []*dto.ViagemDTO{}

	go func() {
		confirm := 0
		for {
			select {
			case resultadoParceialConsulta := <-resultado:
				wg.Done()
				// consultaViagemPlanejamento.ViagensExecutada = append(consultaViagemPlanejamento.ViagensExecutada, resultadoParceialConsulta.ViagensExecutada...)
				consultaViagemPlanejamento.Viagens = append(consultaViagemPlanejamento.Viagens, resultadoParceialConsulta.Viagens...)
				confirm++
				loggerConcorrencia.Debugf("Confirm Ok [%d/%d]", confirm, total)
			case err = <-captura:
				wg.Done()
				confirm++
				loggerConcorrencia.Debugf("Confirm Err [%d/%d]", confirm, total)
			}
		}
	}()
	wg.Wait()

	//TODO - Rever totalização
	consultaViagemPlanejamento.TotExecucao = []int32{int32(len(consultaViagemPlanejamento.Viagens))}

	duracao := time.Since(start)

	var informacoes = make(map[string]interface{})
	informacoes["duracao"] = fmt.Sprintf("%v", duracao)
	consultaViagemPlanejamento.Informacoes = informacoes

	logger.Debugf("QTD Total de Viagens: %d\t em %v\n", len(consultaViagemPlanejamento.Viagens), duracao)
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
	vg.DuracaoSeg = int64(duracao.Seconds())

	if vg.IDHorario > 0 { //Se planejamento encontrado
		diffPartida, diffPartidaFormatada := util.DuracaoEFormatacao(vg.PartidaPlanTime, vgex.Executada.DataInicio)
		vg.DiffPartidaStr = diffPartidaFormatada
		vg.DiffPartida = int64(diffPartida.Seconds())

		diffChegada, diffChegadaFormatada := util.DuracaoEFormatacao(vg.ChegadaPlanTime, vgex.Executada.DataFim)
		vg.DiffChegadaStr = diffChegadaFormatada
		vg.DiffChegada = int64(diffChegada.Seconds())
	}

	vg.VelocidadeMedia = vgex.VelocidadeMedia

	vg.QtdePassageiros = vgex.QntPassageiros

}
