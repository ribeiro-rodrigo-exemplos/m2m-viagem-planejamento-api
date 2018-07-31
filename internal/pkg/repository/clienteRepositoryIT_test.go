package repository

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestCarregarMapaClientes(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestCarregarMapaClientes")

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Reconectar banco de dados devido a falha - %s\n", err)
	}
	c := NewClienteRepository(con)
	clientes, err := c.CarregarMapaClientes()
	if err != nil {
		t.Errorf("Erro ao consultar CarregarMapaClientes - %s\n", err)
	}
	if clientes == nil {
		t.Errorf("Mapa de clientes %v não pode ser nulo\n", clientes)
	}
	if len(clientes) < 1 {
		t.Errorf("Mapa de clientes %v não pode ser vazio\n", clientes)
	}
	for _, cliente := range clientes {
		t.Logf("%+v", cliente)
	}
}
