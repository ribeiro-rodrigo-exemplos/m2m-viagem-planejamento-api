package cache

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestCarregarCacheAgrupamentos(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestCarregarCacheAgrupamentos")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	linhaRepository := repository.NewLinhaRepository(session)
	linhaCache, err := GetLinha(linhaRepository)
	if err != nil {
		t.Errorf("Obter Dependência Cache de Linha - %s\n", err)
	}

	agrupamentoCache, err := GetAgrupamento(linhaCache)
	if err != nil {
		t.Errorf("Obter Cache de Agrupamento - %s\n", err)
	}

	agrupamentos := agrupamentoCache.cache
	if agrupamentos == nil {
		t.Errorf("Cache de agrupamentos %v não pode ser nulo\n", agrupamentos)
	}
	if len(agrupamentos) < 1 {
		t.Errorf("Cache de agrupamentos %v não pode ser vazio\n", agrupamentos)
	}

	if l, _ := agrupamentoCache.Get(38); l == nil {
		t.Errorf("Agrupamentos 38 não encontrado\n")
	}

	a, _ := agrupamentoCache.Get(38)
	if a != nil {
		for _, l := range a.Linhas {
			t.Logf("%+v %+v %+v\n", l.ID, l.Nome, l.Numero)
		}
	}

	for _, agrupamento := range agrupamentos {
		t.Logf("%+v %+v %+v\n", agrupamento.Agrupamento, len(agrupamento.Linhas), len(agrupamento.Trajetos))
	}

	t.Logf("Agrupamentos %d\n", len(agrupamentos))
}

func TestCarregarCacheAgrupamentos39(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestCarregarCacheAgrupamentos")

	qtdTrajetos := 41

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}
	linhaRepository := repository.NewLinhaRepository(session)
	linhaCache, err := GetLinha(linhaRepository)
	if err != nil {
		t.Errorf("Obter Dependência Cache de Linha - %s\n", err)
	}

	agrupamentoCache, err := GetAgrupamento(linhaCache)
	if err != nil {
		t.Errorf("Obter Cache de Agrupamento - %s\n", err)
	}

	agrupamentos := agrupamentoCache.cache
	if agrupamentos == nil {
		t.Errorf("Cache de agrupamentos %v não pode ser nulo\n", agrupamentos)
	}
	if len(agrupamentos) < 1 {
		t.Errorf("Cache de agrupamentos %v não pode ser vazio\n", agrupamentos)
	}

	if l, _ := agrupamentoCache.Get(39); l == nil {
		t.Errorf("Agrupamentos 39 não encontrado\n")
	}

	a, _ := agrupamentoCache.Get(39)

	if a == nil || len(a.MapaTrajetos) != qtdTrajetos {
		t.Errorf("Qtd de trajetos esperada %v, mas encontrada %v\n", qtdTrajetos, len(a.Trajetos))
	}

	if a != nil {
		var idTrajetos map[string]interface{}

		idTrajetos = make(map[string]interface{})
		for _, v := range a.TrajetosDTO {
			idTrajetos[v.ID.Hex()] = v
		}
		if len(idTrajetos) != qtdTrajetos {
			t.Errorf("Qtd de trajetos unicos esperada %v, mas encontrada %v\n", qtdTrajetos, len(idTrajetos))
		}

		idTrajetos = make(map[string]interface{})
		for _, v := range a.Trajetos {
			idTrajetos[v.ID.Hex()] = v
		}
		if len(idTrajetos) != qtdTrajetos {
			t.Errorf("Qtd de trajetos unicos esperada %v, mas encontrada %v\n", qtdTrajetos, len(idTrajetos))
		}

		for _, l := range a.Linhas {
			t.Logf("%+v %+v %+v\n", l.ID, l.Nome, l.Numero)
		}

		t.Logf("Trajetos - %v - %v - %v\n", len(a.MapaTrajetos), len(a.Trajetos), len(a.TrajetosDTO))

		for k, v := range a.MapaTrajetos {
			t.Logf("Mapa TRAJETOS: %v - %v \n", k, v.Nome)
		}

		for _, v := range a.Trajetos {
			t.Logf("ENTITY TRAJETOS: %v - %v \n", v.ID.Hex(), v.Nome)
		}

		for _, v := range a.TrajetosDTO {
			t.Logf("DTO TRAJETOS: %v - %v \n", v.ID.Hex(), v.Descricao)
		}

	}

	for _, agrupamento := range agrupamentos {
		t.Logf("%+v %+v %+v\n", agrupamento.Agrupamento, len(agrupamento.Linhas), len(agrupamento.Trajetos))
	}

	t.Logf("Agrupamentos %d\n", len(agrupamentos))
}
