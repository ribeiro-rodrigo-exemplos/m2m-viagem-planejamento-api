package cache

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"gopkg.in/mgo.v2/bson"
)

func TestConsultarCacheTrajetoVazioECheio(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestConsultarCacheTrajetoVazioECheio")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Obter conexão - %s\n", err)
	}
	trajetoRepository := repository.NewTrajetoRepository(session)

	trajetoCache, err := GetTrajeto(trajetoRepository)
	if err != nil {
		t.Errorf("Obter Cache de Trajeto - %s\n", err)
	}

	var id bson.ObjectId
	var trajeto model.Trajeto

	id = bson.ObjectIdHex("555b6e830850536438063762")
	trajeto, err = trajetoCache.Get(id.Hex())
	if err != nil {
		t.Errorf("Obter entrada de cache - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Cache de trajetos %v não pode ser nulo\n", trajeto)
	}

	id = bson.ObjectIdHex("555b6e830850536438063762")
	trajeto, err = trajetoCache.Get(id.Hex())
	if err != nil {
		t.Errorf("Obter entrada de cache - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Cache de trajetos %v não pode ser nulo\n", trajeto)
	}

	id = bson.ObjectIdHex("555b6e830850536438063761")
	trajeto, err = trajetoCache.Get(id.Hex())
	if err != nil {
		t.Errorf("Obter entrada de cache - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Cache de trajetos %v não pode ser nulo\n", trajeto)
	}

	id = bson.ObjectIdHex("555b6e830850536438063761")
	trajeto, err = trajetoCache.Get(id.Hex())
	if err != nil {
		t.Errorf("Obter entrada de cache - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Cache de trajetos %v não pode ser nulo\n", trajeto)
	}

}

func TestConsultarCacheEPontoFinal(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestConsultarCacheEPontoFinal")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Obter conexão - %s\n", err)
	}
	trajetoRepository := repository.NewTrajetoRepository(session)

	trajetoCache, err := GetTrajeto(trajetoRepository)
	if err != nil {
		t.Errorf("Obter Cache de Trajeto - %s\n", err)
	}

	var id bson.ObjectId
	var idEndPoint bson.ObjectId
	var trajeto model.Trajeto

	id = bson.ObjectIdHex("555b6e830850536438063762")
	idEndPoint = bson.ObjectIdHex("555b32f8085053643806365b")
	trajeto, err = trajetoCache.Get(id.Hex())
	t.Logf("%v\n", trajeto)
	if err != nil {
		t.Errorf("Obter entrada de cache - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Cache de trajetos %v não pode ser nulo\n", trajeto)
	}

	if !trajeto.EndPoint.ID.Valid() {
		t.Errorf("Cache de trajetos com EndPoint inválido %v \n", trajeto.EndPoint.ID)
	}

	if idEndPoint != trajeto.EndPoint.ID {
		t.Errorf("EndPoint esperado %s diferente do obtido %s \n", idEndPoint, trajeto.EndPoint.ID)
	}

	id = bson.ObjectIdHex("555b6e830850536438063761")
	idEndPoint = bson.ObjectIdHex("555b32f8085053643806365f")
	trajeto, err = trajetoCache.Get(id.Hex())
	t.Logf("%v\n", trajeto)
	if err != nil {
		t.Errorf("Obter entrada de cache - %s\n", err)
	}
	if !trajeto.ID.Valid() {
		t.Errorf("Cache de trajetos %v não pode ser nulo\n", trajeto)
	}

	if !trajeto.EndPoint.ID.Valid() {
		t.Errorf("Cache de trajetos com EndPoint inválido %v \n", trajeto.EndPoint.ID)
	}

	if idEndPoint != trajeto.EndPoint.ID {
		t.Errorf("EndPoint esperado %s diferente do obtido %s \n", idEndPoint, trajeto.EndPoint.ID)
	}

}
