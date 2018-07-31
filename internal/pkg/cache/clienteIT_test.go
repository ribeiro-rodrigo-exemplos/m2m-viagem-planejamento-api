package cache

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestCarregarCacheClientes(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestCarregarCacheClientes")

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Obter conexão - %s\n", err)
	}
	clienteRepository := repository.NewClienteRepository(con)

	clienteCache, err := GetCliente(clienteRepository)
	if err != nil {
		t.Errorf("Obter Cache de Cliente - %s\n", err)
	}

	clientes := clienteCache.Cache
	if err != nil {
		t.Errorf("Erro ao consultar CarregarMapaClientes - %s\n", err)
	}
	if clientes == nil {
		t.Errorf("Cache de clientes %v não pode ser nulo\n", clientes)
	}
	if len(clientes) < 1 {
		t.Errorf("Cache de clientes %v não pode ser vazio\n", clientes)
	}
	for _, cliente := range clientes {
		t.Logf("%+v", cliente)
	}
}
