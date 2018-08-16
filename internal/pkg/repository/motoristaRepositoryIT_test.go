package repository

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestCarregarMapaMotoristas(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestCarregarMapaMotoristas")

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Reconectar banco de dados devido a falha - %s\n", err)
	}
	c := NewMotoristaRepository(con)
	motoristas, err := c.CarregarMapaMotoristas()
	if err != nil {
		t.Errorf("Erro ao consultar CarregarMapaMotoristas - %s\n", err)
	}
	if motoristas == nil {
		t.Errorf("Mapa de motoristas %v não pode ser nulo\n", motoristas)
	}
	if len(motoristas) < 1 {
		t.Errorf("Mapa de motoristas %v não pode ser vazio\n", motoristas)
	}
	if len(motoristas) < 100 {
		t.Errorf("Mapa de motoristas não pode ter menos de 100 registros: %v\n", motoristas)
	}
	for _, motorista := range motoristas {
		t.Logf("%+v", motorista)
	}
}
