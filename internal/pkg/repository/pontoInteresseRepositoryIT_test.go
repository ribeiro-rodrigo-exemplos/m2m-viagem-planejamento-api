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

func TestCarregarPontoInteresse(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestCarregarMapaMotoristas")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	c := NewPontoInteresseRepository(session)

	var nomeEsperado string
	var id bson.ObjectId
	var pontoInteresse *model.PontoInteresse

	nomeEsperado = "PU013 - AV. FRANCISCO ERNESTO FÁVERO"
	id = bson.ObjectIdHex("5579d6b2f50beb13664c9cdc")
	pontoInteresse, err = c.ConsultarPorID(id)
	t.Logf("%v\n", pontoInteresse)
	if err != nil {
		t.Errorf("Erro ao consultar pontoInteresse - %s\n", err)
	}
	if pontoInteresse == nil {
		t.Errorf("PontoInteresse não pode ser nulo\n")
	}
	if nomeEsperado != pontoInteresse.Nome {
		t.Errorf("Nome do PontoInteresse esperado %q, mas obtido %q\n", nomeEsperado, pontoInteresse.Nome)
	}

	nomeEsperado = "Terminal Alvorada"
	id = bson.ObjectIdHex("555b32f8085053643806365f")
	pontoInteresse, err = c.ConsultarPorID(id)
	t.Logf("%v\n", pontoInteresse)
	if err != nil {
		t.Errorf("Erro ao consultar pontoInteresse - %s\n", err)
	}
	if pontoInteresse == nil {
		t.Errorf("PontoInteresse não pode ser nulo\n")
	}
	if nomeEsperado != pontoInteresse.Nome {
		t.Errorf("Nome do PontoInteresse esperado %q, mas obtido %q\n", nomeEsperado, pontoInteresse.Nome)
	}

	nomeEsperado = "Terminal Paulo da Portela"
	id = bson.ObjectIdHex("555b32f8085053643806365b")
	pontoInteresse, err = c.ConsultarPorID(id)
	t.Logf("%v\n", pontoInteresse)
	if err != nil {
		t.Errorf("Erro ao consultar pontoInteresse - %s\n", err)
	}
	if pontoInteresse == nil {
		t.Errorf("PontoInteresse não pode ser nulo\n")
	}
	if nomeEsperado != pontoInteresse.Nome {
		t.Errorf("Nome do PontoInteresse esperado %q, mas obtido %q\n", nomeEsperado, pontoInteresse.Nome)
	}

}

func TestCarregarMapaPontoInteresses(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestCarregarMapaPontoInteresses")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	id1 := bson.ObjectIdHex("5579d6b2f50beb13664c9cdc")
	id2 := bson.ObjectIdHex("555b32f8085053643806365f")
	id3 := bson.ObjectIdHex("555b32f8085053643806365b")
	listaIDs := []bson.ObjectId{
		id1,
		id2,
		id3,
	}
	c := NewPontoInteresseRepository(session)
	identificacaoPontosFinal, err := c.CarregarMapaPontoInteresses(listaIDs)
	if err != nil {
		t.Errorf("Erro ao ListarIdentificacaoPontosFinal - %s\n", err)
	}
	if identificacaoPontosFinal == nil {
		t.Errorf("Lista de IdentificacaoPontosFinal não pode ser nula\n")
	}
	if len(identificacaoPontosFinal) < 1 {
		t.Errorf("Lista de IdentificacaoPontosFinal %v não pode ser vazia\n", identificacaoPontosFinal)
	}
	if len(identificacaoPontosFinal) != 3 {
		t.Errorf("Lista de IdentificacaoPontosFinal deve ser 3, mas é %d\n", len(identificacaoPontosFinal))
	}
	// for _, vgex := range identificacaoPontosFinal {
	// 	t.Logf("%+v", vgex)
	// 	fmt.Printf("%s\n", vgex)
	// }

	t.Logf("IDs Ponto Final %d\n", len(identificacaoPontosFinal))
}

func TestListarIdentificacaoPontosFinal(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestListarIdentificacaoPontosFinal")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}

	c := NewPontoInteresseRepository(session)
	identificacaoPontosFinal, err := c.ListarIdentificacaoPontosFinal()
	if err != nil {
		t.Errorf("Erro ao ListarIdentificacaoPontosFinal - %s\n", err)
	}
	if identificacaoPontosFinal == nil {
		t.Errorf("Lista de IdentificacaoPontosFinal não pode ser nula\n")
	}
	if len(identificacaoPontosFinal) < 1 {
		t.Errorf("Lista de IdentificacaoPontosFinal %v não pode ser vazia\n", identificacaoPontosFinal)
	}
	// for _, vgex := range identificacaoPontosFinal {
	// 	t.Logf("%+v", vgex)
	// 	fmt.Printf("%s\n", vgex)
	// }

	t.Logf("IDs Ponto Final %d\n", len(identificacaoPontosFinal))
}
