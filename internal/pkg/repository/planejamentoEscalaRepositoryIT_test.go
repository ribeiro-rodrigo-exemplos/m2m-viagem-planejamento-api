package repository

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"gopkg.in/mgo.v2/bson"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestListarPlanejamentosEscalas(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestListarPlanejamentosEscalas")

	filter := &dto.FilterDTO{
		ListaTrajetos: []bson.ObjectId{
			bson.ObjectIdHex("555b6e830850536438063762"),
			// bson.ObjectIdHex("555b6e830850536438063761"),
		},
		IDCliente:  209,
		IDVeiculo:  150,
		Ordenacao:  []string{"veiculo", "data"},
		DataInicio: "2018-07-24 18:00:00",
		DataFim:    "2018-07-24 23:59:59",
		TipoDia:    []string{"O", "E", "3", "U"},
	}

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Reconectar banco de dados devido a falha - %s\n", err)
	}
	c := NewPlanejamentoEscalaRepository(con)
	planejamentosEscala, err := c.ListarPlanejamentosEscala(filter)
	if err != nil {
		t.Errorf("Erro ao ListarPlanejamentosEscala - %s\n", err)
	}
	if planejamentosEscala == nil {
		t.Errorf("Lista de PlanejamentoEscala %v não pode ser nula\n", planejamentosEscala)
	}
	if len(planejamentosEscala) < 1 {
		t.Errorf("Lista de PlanejamentoEscala %v não pode ser  vazia\n", planejamentosEscala)
	}
	for _, cliente := range planejamentosEscala {
		t.Logf("%+v", cliente)
	}
}