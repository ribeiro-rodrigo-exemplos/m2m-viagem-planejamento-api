package cache

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestListarCacheLinhas(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestListarCacheLinhas")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	linhaRepository := repository.NewLinhaRepository(session)

	linhaCache, err := GetLinha(linhaRepository)
	if err != nil {
		t.Errorf("Obter Cache de Linha - %s\n", err)
	}

	linhas, _ := linhaCache.ListAll()
	if err != nil {
		t.Errorf("Obter cache de lista de Linhas - %s\n", err)
	}

	if linhas == nil {
		t.Errorf("Cache de linhas %v não pode ser nulo\n", linhas)
	}
	if len(linhas) < 1 {
		t.Errorf("Cache de linhas %v não pode ser vazio\n", linhas)
	}

	// for _, linha := range linhas {
	// 	t.Logf("%+v", linha)
	// }

	t.Logf("Linhas %d\n", len(linhas))
}

func TestCarregarCacheLinhas(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestCarregarCacheLinhas")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	linhaRepository := repository.NewLinhaRepository(session)

	linhaCache, err := GetLinha(linhaRepository)
	if err != nil {
		t.Errorf("Obter Cache de Linha - %s\n", err)
	}

	linhas := linhaCache.cache
	if linhas == nil {
		t.Errorf("Cache de linhas %v não pode ser nulo\n", linhas)
	}
	if len(linhas) < 1 {
		t.Errorf("Cache de linhas %v não pode ser vazio\n", linhas)
	}

	if l, _ := linhaCache.Get(bson.ObjectIdHex("555b6e830850536438063763")); !l.ID.Valid() {
		t.Errorf("Linhas 555b6e830850536438063763 não encontrada\n")
	}
	a, _ := linhaCache.Get(bson.ObjectIdHex("555b6e830850536438063763"))
	if a.ID.Valid() {
		t.Logf("%+v - %+v - %+v\n", a.Agrupamento, a.ID, a.Nome)
	}

	// for _, linha := range linhas {
	// 	t.Logf("%+v", linha)
	// }

	t.Logf("Linhas %d\n", len(linhas))
}
