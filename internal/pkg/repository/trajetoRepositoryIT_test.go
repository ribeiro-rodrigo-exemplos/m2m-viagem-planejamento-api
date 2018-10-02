package repository

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"gopkg.in/mgo.v2/bson"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestCarregarTrajeto(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestCarregarMapaMotoristas")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	c := NewTrajetoRepository(session)

	var nomeEsperado string
	var id bson.ObjectId
	var trajeto model.Trajeto

	nomeEsperado = "35 - ALVORADA  X  MADUREIRA ( PARADOR ) - Volta"
	id = bson.ObjectIdHex("555b6e830850536438063762")
	trajeto, err = c.ConsultarPorID(id)
	t.Logf("%v\n", trajeto)
	if err != nil {
		t.Errorf("Erro ao consultar trajeto - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Trajeto não pode ser nulo\n")
	}
	if nomeEsperado != trajeto.Nome {
		t.Errorf("Nome do Trajeto esperado %q, mas obtido %q\n", nomeEsperado, trajeto.Nome)
	}

	nomeEsperado = "35 - MADUREIRA X ALVORADA ( PARADOR ) - ida"
	id = bson.ObjectIdHex("555b6e830850536438063761")
	trajeto, err = c.ConsultarPorID(id)
	t.Logf("%v\n", trajeto)
	if err != nil {
		t.Errorf("Erro ao consultar trajeto - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Trajeto não pode ser nulo\n")
	}
	if nomeEsperado != trajeto.Nome {
		t.Errorf("Nome do Trajeto esperado %q, mas obtido %q\n", nomeEsperado, trajeto.Nome)
	}

	nomeEsperado = "10- SANTA CRUZ X ALVORADA ( EXPRESSO )"
	id = bson.ObjectIdHex("555b4eadaecc1a6638f3ab29")
	trajeto, err = c.ConsultarPorID(id)
	t.Logf("%v\n", trajeto)
	if err != nil {
		t.Errorf("Erro ao consultar trajeto - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Trajeto não pode ser nulo\n")
	}
	if nomeEsperado != trajeto.Nome {
		t.Errorf("Nome do Trajeto esperado %q, mas obtido %q\n", nomeEsperado, trajeto.Nome)
	}

}

func TestCarregarMapaTrajetos(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestCarregarMapaTrajetos")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}

	var nomeEsperado string
	var id bson.ObjectId
	var trajeto model.Trajeto

	c := NewTrajetoRepository(session)
	mapaTrajetos, err := c.CarregarMapaTrajetos()
	if err != nil {
		t.Errorf("Erro ao ListarIdentificacaoPontosFinal - %s\n", err)
	}
	if mapaTrajetos == nil {
		t.Errorf("Mapa de Trajetos não pode ser nulo\n")
	}
	if len(mapaTrajetos) < 1 {
		t.Errorf("Mapa de Trajetos %v não pode estar vazio\n", mapaTrajetos)
	}

	nomeEsperado = "35 - ALVORADA  X  MADUREIRA ( PARADOR ) - Volta"
	id = bson.ObjectIdHex("555b6e830850536438063762")
	trajeto = mapaTrajetos[id.Hex()]
	if nomeEsperado != trajeto.Nome {
		t.Errorf("Nome do Trajeto esperado %q, mas obtido %q - %+v\n", nomeEsperado, trajeto.Nome, trajeto)
	}

	nomeEsperado = "35 - MADUREIRA X ALVORADA ( PARADOR ) - ida"
	id = bson.ObjectIdHex("555b6e830850536438063761")
	trajeto = mapaTrajetos[id.Hex()]
	if nomeEsperado != trajeto.Nome {
		t.Errorf("Nome do Trajeto esperado %q, mas obtido %q - %+v\n", nomeEsperado, trajeto.Nome, trajeto)
	}

	// for _, vgex := range identificacaoPontosFinal {
	// 	t.Logf("%+v", vgex)
	// 	fmt.Printf("%s\n", vgex)
	// }

	t.Logf("Trajetos %d\n", len(mapaTrajetos))
}

func TestListarTrajetos(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestListarTrajetos")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}

	c := NewTrajetoRepository(session)
	listaTrajetos, err := c.ListarTrajetos()
	if err != nil {
		t.Errorf("Erro ao ListarTrajetos - %s\n", err)
	}
	if listaTrajetos == nil {
		t.Errorf("Lista de Trajetos não pode ser nula\n")
	}
	if len(listaTrajetos) < 1 {
		t.Errorf("Lista de Trajetos %v não pode ser vazia\n", listaTrajetos)
	}
	// for _, vgex := range identificacaoPontosFinal {
	// 	t.Logf("%+v", vgex)
	// 	fmt.Printf("%s\n", vgex)
	// }

	t.Logf("Trajetos %d\n", len(listaTrajetos))
}
