package viagemplanejamento

import (
	"sync"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/cache"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"gopkg.in/mgo.v2/bson"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestConsultarViagemPlanejamentoPorUmTrajeto(t *testing.T) {
	cfg.InitConfig("../../../../configs/config.json")
	logger = logging.NewLogger("service.viagemplanejamento", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	cache.InitConfig()
	t.Log("TestConsultarViagemPlanejamentoPorUmTrajeto")

	var err error

	cacheCliente, _ := cache.GetCliente(nil)
	cliente := cacheCliente.Cache[209]

	cacheMotorista, _ := cache.GetMotorista(nil)
	cacheTrajeto, _ := cache.GetTrajeto(nil)
	cachePontoInteresse, _ := cache.GetPontoInteresse(nil)
	cacheAgrupamento, _ := cache.GetAgrupamento(nil)

	id := bson.ObjectIdHex("555b6e830850536438063762")
	dataInicio := "2018-07-24 18:00:00"
	dataFim := "2018-07-24 23:59:59"
	filter := dto.FilterDTO{
		ListaTrajetos: []dto.TrajetoDTO{
			dto.TrajetoDTO{ID: &id},
			// bson.ObjectIdHex("555b6e830850536438063761"),
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  "horario",
		DataInicio: &dataInicio,
		DataFim:    &dataFim,
		Complemento: dto.DadosComplementares{
			Cliente: cliente,
		},
	}
	filter.TipoDia = model.TiposDia.FromDate(filter.GetDataInicio(), []string{"O", "F"})

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	planejamentoEscalaRepository := repository.NewPlanejamentoEscalaRepository(con)

	mongoDB, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	viagemExecutadaRepository := repository.NewViagemExecutadaRepository(mongoDB)

	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente, cacheMotorista, cacheTrajeto, cachePontoInteresse, cacheAgrupamento)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	resultado := make(chan *dto.ConsultaViagemPlanejamentoDTO, 5)
	captura := make(chan error, 5)
	var wg sync.WaitGroup
	// for index := 0; index < 2; index++ {
	// 	wg.Add(2)
	// 	go vps.ConsultarPorTrajeto(filter, resultado, captura)
	// 	go vps.ConsultarPorTrajeto(filter, resultado, captura)
	// }
	for index := 0; index < 1; index++ {
		wg.Add(1)
		go vps.ConsultarPorTrajeto(filter, resultado, captura)
	}
	go func() {
		confirm := 0
		for {
			select {
			case consultaViagemPlanejamento = <-resultado:
				wg.Done()
				confirm++
				t.Log("Confirm Ok", confirm)
			case err = <-captura:
				wg.Done()
				confirm++
				t.Log("Confirm Err", confirm)
			}
		}
	}()
	wg.Wait()

	if err != nil {
		t.Errorf("Erro ao ConsultarViagemPlanejamento - %s\n", err)
	}
	if consultaViagemPlanejamento == nil {
		t.Errorf("Consulta de ViagemPlanejamento não pode ser nula\n")
		return
	}
	if consultaViagemPlanejamento.Viagens == nil {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser nula\n", consultaViagemPlanejamento.Viagens)
	}
	if len(consultaViagemPlanejamento.Viagens) < 1 {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser vazia\n", consultaViagemPlanejamento.Viagens)
	}
	for _, vg := range consultaViagemPlanejamento.Viagens {
		t.Logf("%+v\n", vg)
	}
}

func TestConsultarViagemPlanejamentoPorUmTrajetoEmUmaNoite(t *testing.T) {
	cfg.InitConfig("../../../../configs/config.json")
	InitConfig()
	database.InitConfig()

	repository.InitConfig()
	cache.InitConfig()
	t.Log("TestConsultarViagemPlanejamentoPorUmTrajetoEmUmaNoite")

	id := bson.ObjectIdHex("555b6e830850536438063762")
	dataInicio := "2018-07-24 18:00:00"
	dataFim := "2018-07-24 23:59:59"
	var err error
	filter := dto.FilterDTO{
		ListaTrajetos: []dto.TrajetoDTO{
			dto.TrajetoDTO{
				ID:    &id,
				Linha: dto.LinhaDTO{Numero: "5702A1"},
			},
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  "horario",
		DataInicio: &dataInicio,
		DataFim:    &dataFim,
	}

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	planejamentoEscalaRepository := repository.NewPlanejamentoEscalaRepository(con)

	mongoDB, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	viagemExecutadaRepository := repository.NewViagemExecutadaRepository(mongoDB)

	cacheCliente, _ := cache.GetCliente(nil)
	cacheMotorista, _ := cache.GetMotorista(nil)
	cacheTrajeto, _ := cache.GetTrajeto(nil)
	cachePontoInteresse, _ := cache.GetPontoInteresse(nil)
	cacheAgrupamento, _ := cache.GetAgrupamento(nil)

	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente, cacheMotorista, cacheTrajeto, cachePontoInteresse, cacheAgrupamento)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	consultaViagemPlanejamento, err = vps.ConsultarPeriodo(filter)

	if err != nil {
		t.Errorf("Erro ao ConsultarViagemPlanejamento - %s\n", err)
	}
	if consultaViagemPlanejamento == nil {
		t.Errorf("Consulta de ViagemPlanejamento não pode ser nula\n")
		return
	}
	if consultaViagemPlanejamento.Viagens == nil {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser nula\n", consultaViagemPlanejamento.Viagens)
	}
	if len(consultaViagemPlanejamento.Viagens) < 1 {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser vazia\n", consultaViagemPlanejamento.Viagens)
	}
	if consultaViagemPlanejamento.Totalizadores.Canceladas != 1 {
		t.Errorf("Totalizador canceladas não pode ser diferente de 1, mas foi %v \n", consultaViagemPlanejamento.Totalizadores.Canceladas)
	}

	for _, vg := range consultaViagemPlanejamento.Viagens {
		t.Logf("%+v\n", vg)
	}
}

func TestConsultarViagemPlanejamentoPorUmAgrupamentoComEmpresa(t *testing.T) {
	cfg.InitConfig("../../../../configs/config.json")
	InitConfig()
	database.InitConfig()

	repository.InitConfig()
	cache.InitConfig()
	t.Log("TestConsultarViagemPlanejamentoPorUmAgrupamento")

	agrupamento := int32(38)
	id := bson.ObjectIdHex("555b6e830850536438063762")
	// dataInicio := "2018-08-24 18:00:00"
	dataInicio := "2018-08-24 00:00:00"
	dataFim := "2018-08-24 23:59:59"
	var err error
	filter := dto.FilterDTO{

		ListaAgrupamentos: []dto.AgrupamentoDTO{
			dto.AgrupamentoDTO{
				ID: agrupamento,
			},
		},
		ListaTrajetos: []dto.TrajetoDTO{
			dto.TrajetoDTO{
				ID:    &id,
				Linha: dto.LinhaDTO{Numero: "5702A1"}, // seleção de agrupamento deve olhar p cache de linha
			},
		},
		ListaEmpresas: []dto.EmpresaDTO{
			dto.EmpresaDTO{
				ID: 1851,
			},
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  "horario",
		DataInicio: &dataInicio,
		DataFim:    &dataFim,
	}

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	planejamentoEscalaRepository := repository.NewPlanejamentoEscalaRepository(con)

	mongoDB, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	viagemExecutadaRepository := repository.NewViagemExecutadaRepository(mongoDB)

	cacheCliente, _ := cache.GetCliente(nil)
	cacheMotorista, _ := cache.GetMotorista(nil)
	cacheTrajeto, _ := cache.GetTrajeto(nil)
	cachePontoInteresse, _ := cache.GetPontoInteresse(nil)
	cacheAgrupamento, _ := cache.GetAgrupamento(nil)

	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente, cacheMotorista, cacheTrajeto, cachePontoInteresse, cacheAgrupamento)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	consultaViagemPlanejamento, err = vps.ConsultarPeriodo(filter)

	if err != nil {
		t.Errorf("Erro ao ConsultarViagemPlanejamento - %s\n", err)
	}
	if consultaViagemPlanejamento == nil {
		t.Errorf("Consulta de ViagemPlanejamento não pode ser nula\n")
		return
	}
	if consultaViagemPlanejamento.Viagens == nil {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser nula\n", consultaViagemPlanejamento.Viagens)
	}
	if len(consultaViagemPlanejamento.Viagens) < 1 {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser vazia\n", consultaViagemPlanejamento.Viagens)
	}

	for _, vg := range consultaViagemPlanejamento.Viagens {
		t.Logf("%+v\n", vg)
	}

	t.Logf("%d\n", len(consultaViagemPlanejamento.Viagens))
}

func TestConsultarViagemPlanejamentoPorAgrupamentoInexistente(t *testing.T) {
	cfg.InitConfig("../../../../configs/config.json")
	InitConfig()
	database.InitConfig()

	repository.InitConfig()
	cache.InitConfig()
	t.Log("TestConsultarViagemPlanejamentoPorUmAgrupamento")

	agrupamento := int32(999999999)
	dataInicio := "2018-08-24 00:00:00"
	dataFim := "2018-08-24 23:59:59"
	var err error
	filter := dto.FilterDTO{

		ListaAgrupamentos: []dto.AgrupamentoDTO{
			dto.AgrupamentoDTO{
				ID: agrupamento,
			},
		},
		ListaEmpresas: []dto.EmpresaDTO{
			dto.EmpresaDTO{
				ID: 1851,
			},
		},
		IDCliente:  209,
		Ordenacao:  "horario",
		DataInicio: &dataInicio,
		DataFim:    &dataFim,
	}

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	planejamentoEscalaRepository := repository.NewPlanejamentoEscalaRepository(con)

	mongoDB, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	viagemExecutadaRepository := repository.NewViagemExecutadaRepository(mongoDB)

	cacheCliente, _ := cache.GetCliente(nil)
	cacheMotorista, _ := cache.GetMotorista(nil)
	cacheTrajeto, _ := cache.GetTrajeto(nil)
	cachePontoInteresse, _ := cache.GetPontoInteresse(nil)
	cacheAgrupamento, _ := cache.GetAgrupamento(nil)

	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente, cacheMotorista, cacheTrajeto, cachePontoInteresse, cacheAgrupamento)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	consultaViagemPlanejamento, err = vps.ConsultarPeriodo(filter)

	if err != nil {
		t.Errorf("Erro ao ConsultarViagemPlanejamento - %s\n", err)
	}
	if consultaViagemPlanejamento == nil {
		t.Errorf("Consulta de ViagemPlanejamento não pode ser nula\n")
		return
	}
	if len(consultaViagemPlanejamento.Viagens) > 0 {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v deve ser vazia\n", consultaViagemPlanejamento.Viagens)
	}

	t.Logf("%d\n", len(consultaViagemPlanejamento.Viagens))
}

func TestConsultarViagemPlanejamentoPorDoisTrajetosEmUmDia(t *testing.T) {
	cfg.InitConfig("../../../../configs/config.json")
	InitConfig()
	database.InitConfig()
	repository.InitConfig()
	cache.InitConfig()
	t.Log("TestConsultarViagemPlanejamentoPorUmTrajetoEmUmDia")

	var err error

	id1 := bson.ObjectIdHex("555b6e830850536438063762")
	id2 := bson.ObjectIdHex("555b6e830850536438063761")
	dataInicio := "2018-08-02 00:00:00"
	dataFim := "2018-08-02 23:59:59"
	filter := dto.FilterDTO{
		ListaTrajetos: []dto.TrajetoDTO{
			dto.TrajetoDTO{ID: &id1},
			dto.TrajetoDTO{ID: &id2},
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  "horario",
		DataInicio: &dataInicio,
		DataFim:    &dataFim,
	}

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	planejamentoEscalaRepository := repository.NewPlanejamentoEscalaRepository(con)

	mongoDB, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	viagemExecutadaRepository := repository.NewViagemExecutadaRepository(mongoDB)

	cacheCliente, _ := cache.GetCliente(nil)
	cacheMotorista, _ := cache.GetMotorista(nil)
	cacheTrajeto, _ := cache.GetTrajeto(nil)
	cachePontoInteresse, _ := cache.GetPontoInteresse(nil)
	cacheAgrupamento, _ := cache.GetAgrupamento(nil)

	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente, cacheMotorista, cacheTrajeto, cachePontoInteresse, cacheAgrupamento)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	consultaViagemPlanejamento, err = vps.ConsultarPeriodo(filter)

	if err != nil {
		t.Errorf("Erro ao ConsultarViagemPlanejamento - %s\n", err)
	}
	if consultaViagemPlanejamento == nil {
		t.Errorf("Consulta de ViagemPlanejamento não pode ser nula\n")
		return
	}
	if consultaViagemPlanejamento.Viagens == nil {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser nula\n", consultaViagemPlanejamento.Viagens)
	}
	if len(consultaViagemPlanejamento.Viagens) < 1 {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser vazia\n", consultaViagemPlanejamento.Viagens)
	}
	if consultaViagemPlanejamento.Totalizadores.Canceladas != 9 {
		t.Errorf("Totalizador canceladas não pode ser diferente de 9, mas foi %v \n", consultaViagemPlanejamento.Totalizadores.Canceladas)
	}

	for _, vg := range consultaViagemPlanejamento.Viagens {
		t.Logf("%+v\n", vg)
	}
}

func TestConsultarViagemPlanejamentoPorUmTrajetoEmSeteDias(t *testing.T) {
	cfg.InitConfig("../../../../configs/config.json")
	InitConfig()
	database.InitConfig()
	repository.InitConfig()
	cache.InitConfig()
	t.Log("TestConsultarViagemPlanejamentoPorUmTrajetoEmSeteDias")

	var err error

	id := bson.ObjectIdHex("555b6e830850536438063762")
	dataInicio := "2018-07-24 18:00:00"
	dataFim := "2018-08-02 20:00:00"
	filter := dto.FilterDTO{
		ListaTrajetos: []dto.TrajetoDTO{
			dto.TrajetoDTO{ID: &id},
			// bson.ObjectIdHex("555b6e830850536438063761"),
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  "horario",
		DataInicio: &dataInicio,
		DataFim:    &dataFim,
		// DataFim:    "2018-07-24 23:59:59",
		// DataFim: "2018-07-25 17:59:59",
	}

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	planejamentoEscalaRepository := repository.NewPlanejamentoEscalaRepository(con)

	mongoDB, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	viagemExecutadaRepository := repository.NewViagemExecutadaRepository(mongoDB)

	cacheCliente, _ := cache.GetCliente(nil)
	cacheMotorista, _ := cache.GetMotorista(nil)
	cacheTrajeto, _ := cache.GetTrajeto(nil)
	cachePontoInteresse, _ := cache.GetPontoInteresse(nil)
	cacheAgrupamento, _ := cache.GetAgrupamento(nil)

	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente, cacheMotorista, cacheTrajeto, cachePontoInteresse, cacheAgrupamento)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	consultaViagemPlanejamento, err = vps.ConsultarPeriodo(filter)

	if err != nil {
		t.Errorf("Erro ao ConsultarViagemPlanejamento - %s\n", err)
	}
	if consultaViagemPlanejamento == nil {
		t.Errorf("Consulta de ViagemPlanejamento não pode ser nula\n")
		return
	}
	if consultaViagemPlanejamento.Viagens == nil {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser nula\n", consultaViagemPlanejamento.Viagens)
	}
	if len(consultaViagemPlanejamento.Viagens) < 1 {
		t.Errorf("Viagens de Consulta de ViagemPlanejamento %v não pode ser vazia\n", consultaViagemPlanejamento.Viagens)
	}
	for _, vg := range consultaViagemPlanejamento.Viagens {
		t.Logf("%+v\n", vg)
	}
}
