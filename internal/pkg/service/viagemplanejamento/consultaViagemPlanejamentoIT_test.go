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

	filter := dto.FilterDTO{
		ListaTrajetos: []bson.ObjectId{
			bson.ObjectIdHex("555b6e830850536438063762"),
			// bson.ObjectIdHex("555b6e830850536438063761"),
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  []string{"veiculo", "data"},
		DataInicio: "2018-07-24 18:00:00",
		DataFim:    "2018-07-24 23:59:59",
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

	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente)

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

	var err error
	filter := dto.FilterDTO{
		ListaTrajetos: []bson.ObjectId{
			bson.ObjectIdHex("555b6e830850536438063762"),
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  []string{"veiculo", "data"},
		DataInicio: "2018-07-24 18:00:00",
		DataFim:    "2018-07-24 23:59:59",
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
	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	consultaViagemPlanejamento, err = vps.Consultar(filter)

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
		t.Errorf("Totalizador canceladas não pode ser 0\n")
	}

	for _, vg := range consultaViagemPlanejamento.Viagens {
		t.Logf("%+v\n", vg)
	}
}

func TestConsultarViagemPlanejamentoPorDoisTrajetosEmUmDia(t *testing.T) {
	cfg.InitConfig("../../../../configs/config.json")
	InitConfig()
	database.InitConfig()
	repository.InitConfig()
	cache.InitConfig()
	t.Log("TestConsultarViagemPlanejamentoPorUmTrajetoEmUmDia")

	var err error

	filter := dto.FilterDTO{
		ListaTrajetos: []bson.ObjectId{
			bson.ObjectIdHex("555b6e830850536438063762"),
			bson.ObjectIdHex("555b6e830850536438063761"),
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  []string{"veiculo", "data"},
		DataInicio: "2018-08-02 00:00:00",
		DataFim:    "2018-08-02 23:59:59",
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
	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	consultaViagemPlanejamento, err = vps.Consultar(filter)

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
	if consultaViagemPlanejamento.Totalizadores.Canceladas < 1 {
		t.Errorf("Totalizador canceladas não pode ser 0\n")
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

	filter := dto.FilterDTO{
		ListaTrajetos: []bson.ObjectId{
			bson.ObjectIdHex("555b6e830850536438063762"),
			// bson.ObjectIdHex("555b6e830850536438063761"),
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  []string{"veiculo", "data"},
		DataInicio: "2018-07-24 18:00:00",
		// DataFim:    "2018-07-24 23:59:59",
		// DataFim: "2018-07-25 17:59:59",
		DataFim: "2018-08-02 20:00:00",
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
	vps := NewViagemPlanejamentoService(planejamentoEscalaRepository, viagemExecutadaRepository, cacheCliente)

	var consultaViagemPlanejamento *dto.ConsultaViagemPlanejamentoDTO

	consultaViagemPlanejamento, err = vps.Consultar(filter)

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
