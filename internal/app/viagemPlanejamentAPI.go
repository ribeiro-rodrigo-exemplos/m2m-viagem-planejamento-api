package app

import (
	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/cache"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/intercept"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/service/viagemplanejamento"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/webservice"
)

var logger logging.Logger

//Bootstrap - é responsável por iniciar a aplicação
func Bootstrap() {
	logger = logging.NewLogger("", cfg.Config.Logging.Level)

	logger.Infof("Iniciando a aplicação...\n")

	// menssageria.ConectaRabbitmq()
	database.InitConfig()
	repository.InitConfig()
	cache.InitConfig()
	viagemplanejamento.InitConfig()
	intercept.InitConfig()
	webservice.InitConfig()
	webservice.InitServer()
}
