package intercept

import (
	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

var logger logging.Logger
var loggerRequestBody logging.Logger

//InitConfig - é responsável por iniciar configuração da package
func InitConfig() {
	logger = logging.NewLogger("intercept", cfg.Config.Logging.Level)
	loggerRequestBody = logging.NewLogger("intercept.REQUEST_BODY", cfg.Config.Logging.Level)
}
