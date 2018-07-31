package repository

import (
	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

var logger logging.Logger

//InitConfig - é responsável por iniciar configuração da package
func InitConfig() {
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
}
