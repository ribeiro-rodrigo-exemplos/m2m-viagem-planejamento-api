package cache

import (
	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
)

var logger logging.Logger

//InitConfig - é responsável por iniciar configuração da package
func InitConfig() {
	logger = logging.NewLogger("cache", cfg.Config.Logging.Level)
	// client = configuraClient()

	//TODO - Receber dependência conexão como parâmetro
	con, err := database.GetSQLConnection()
	if err != nil {
		logger.Errorf("Obter conexão - %s\n", err)
	}
	clienteRepository := repository.NewClienteRepository(con)

	_, err = GetCliente(clienteRepository)
	if err != nil {
		logger.Errorf("Obter Cache de Cliente - %s\n", err)
	}

	//TODO - Receber dependência conexão como parâmetro
	con, err = database.GetSQLConnection()
	if err != nil {
		logger.Errorf("Obter conexão - %s\n", err)
	}
	motoristaRepository := repository.NewMotoristaRepository(con)

	_, err = GetMotorista(motoristaRepository)
	if err != nil {
		logger.Errorf("Obter Cache de Cliente - %s\n", err)
	}

	//TODO - Receber dependência conexão como parâmetro
	session, err := database.GetMongoDB()
	if err != nil {
		logger.Errorf("Obter conexão - %s\n", err)
	}
	trajetoRepository := repository.NewTrajetoRepository(session)

	_, err = GetTrajeto(trajetoRepository)
	if err != nil {
		logger.Errorf("Obter Cache de Trajeto - %s\n", err)
	}

	pontoInteresseRepository := repository.NewPontoInteresseRepository(session)

	_, err = GetPontoInteresse(pontoInteresseRepository)
	if err != nil {
		logger.Errorf("Obter Cache de Ponto de Interesse - %s\n", err)
	}
}
