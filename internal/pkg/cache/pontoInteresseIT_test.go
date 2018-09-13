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

func TestMapKeyReferenc(t *testing.T) {

	m := make(map[bson.ObjectId]string)

	id1 := bson.ObjectIdHex("5579d6b2f50beb13664c9cdc")
	id2 := bson.ObjectIdHex("5579d6b2f50beb13664c9cdc")

	k1 := id1
	k2 := id2

	ori := "t1"
	exp := "t2"

	m[k1] = ori
	m[k2] = exp

	if len(m) != 1 {
		t.Errorf("Mapa deveria ter 1, mas possui %d elementos", len(m))
	}

	if m[k1] != exp {
		t.Errorf("Valor deveria ser %q, mas é %q", exp, m[k1])
	}

}

func TestConsultarCacheVazioECheio(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestConsultarCacheVazioECheio")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Obter conexão - %s\n", err)
	}
	pontoInteresseRepository := repository.NewPontoInteresseRepository(session)

	pontoInteresseCache, err := GetPontoInteresse(pontoInteresseRepository)
	if err != nil {
		t.Errorf("Obter Cache de PontoInteresse - %s\n", err)
	}

	var id bson.ObjectId
	var pontoInteresse *model.PontoInteresse

	id = bson.ObjectIdHex("555b32f8085053643806365f")
	pontoInteresse, err = pontoInteresseCache.Get(id)
	fmt.Printf("%+v\n", pontoInteresse)
	if err != nil {
		t.Errorf("Obter entrada de chace - %s\n", err)
	}
	if pontoInteresse == nil {
		t.Errorf("Cache de pontoInteresses %v não pode ser nulo\n", pontoInteresse)
	}

	id = bson.ObjectIdHex("555b32f8085053643806365f")
	pontoInteresse, err = pontoInteresseCache.Get(id)
	fmt.Printf("%+v\n", pontoInteresse)
	if err != nil {
		t.Errorf("Obter entrada de chace - %s\n", err)
	}
	if pontoInteresse == nil {
		t.Errorf("Cache de pontoInteresses %v não pode ser nulo\n", pontoInteresse)
	}

	id = bson.ObjectIdHex("5579d6b2f50beb13664c9cdc")
	pontoInteresse, err = pontoInteresseCache.Get(id)
	fmt.Printf("%+v\n", pontoInteresse)
	if err != nil {
		t.Errorf("Obter entrada de chace - %s\n", err)
	}
	if pontoInteresse == nil {
		t.Errorf("Cache de pontoInteresses %v não pode ser nulo\n", pontoInteresse)
	}

	id = bson.ObjectIdHex("5579d6b2f50beb13664c9cdc")
	pontoInteresse, err = pontoInteresseCache.Get(id)
	fmt.Printf("%+v\n", pontoInteresse)
	if err != nil {
		t.Errorf("Obter entrada de chace - %s\n", err)
	}
	if pontoInteresse == nil {
		t.Errorf("Cache de pontoInteresses %v não pode ser nulo\n", pontoInteresse)
	}

}
