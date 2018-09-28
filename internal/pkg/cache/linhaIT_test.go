package cache

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

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
	motoristaRepository := repository.NewLinhaRepository(session)

	motoristaCache, err := GetLinha(motoristaRepository)
	if err != nil {
		t.Errorf("Obter Cache de Linha - %s\n", err)
	}

	motoristas := motoristaCache.Cache
	if err != nil {
		t.Errorf("Erro ao consultar CarregarMapaLinhas - %s\n", err)
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
