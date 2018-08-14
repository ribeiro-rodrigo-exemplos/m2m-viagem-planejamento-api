package repository

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"gopkg.in/mgo.v2/bson"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestListarViagensPorTrajetoUmDia(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestListarViagensPorTrajetoUmDia")

	cliente := &model.Cliente{
		IDCliente: 209,
		Nome:      "BRT",
		Timezone:  "",
	}
	cliente.AtualizarLocation()

	filter := dto.FilterDTO{
		ListaTrajetos: []dto.TrajetoDTO{
			dto.TrajetoDTO{ID: bson.ObjectIdHex("555b6e830850536438063762")},
			// bson.ObjectIdHex("555b6e830850536438063761"),
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  []string{"veiculo", "data"},
		DataInicio: "2018-07-24 18:00:00",
		DataFim:    "2018-07-24 23:59:59",
		TipoDia:    []string{"O", "E", "3", "U"},
		Complemento: dto.DadosComplementares{
			Cliente: cliente,
		},
	}

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}

	viagemExecutadaRepository := NewViagemExecutadaRepository(session)
	viagensExecutada, err := viagemExecutadaRepository.ListarViagensPor(filter)
	if err != nil {
		t.Errorf("Erro ao ListarViagensPor %+v - %s\n", filter, err)
	}
	if viagensExecutada == nil {
		t.Errorf("Lista de ViagemExecutada %v não pode ser nula\n", viagensExecutada)
	}
	if len(viagensExecutada) < 1 {
		t.Errorf("Lista de ViagemExecutada %v não pode ser vazia\n", viagensExecutada)
	}
	for _, vgex := range viagensExecutada {
		t.Logf("%+v", vgex)
	}
}
