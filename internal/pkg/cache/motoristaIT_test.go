package cache

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestCarregarCacheMotoristas(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	repository.InitConfig()
	fmt.Println("TestCarregarCacheMotoristas")

	con, err := database.GetSQLConnection()
	if err != nil {
		t.Errorf("Obter conexão - %s\n", err)
	}
	motoristaRepository := repository.NewMotoristaRepository(con)

	motoristaCache, err := GetMotorista(motoristaRepository)
	if err != nil {
		t.Errorf("Obter Cache de Motorista - %s\n", err)
	}

	motoristas := motoristaCache.Cache
	if err != nil {
		t.Errorf("Erro ao consultar CarregarMapaMotoristas - %s\n", err)
	}
	if motoristas == nil {
		t.Errorf("Cache de motoristas %v não pode ser nulo\n", motoristas)
	}
	if len(motoristas) < 1 {
		t.Errorf("Cache de motoristas %v não pode ser vazio\n", motoristas)
	}
	for _, motorista := range motoristas {
		t.Logf("%+v", motorista)
	}
}
