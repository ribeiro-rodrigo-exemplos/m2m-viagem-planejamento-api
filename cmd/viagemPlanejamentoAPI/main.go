package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/app"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const productionEnv = "PRODUCTION"

func init() {
	http.DefaultClient.Timeout = time.Minute * 2

	cfg.InitConfig("")

	if cfg.Config.RuntimeEnv == productionEnv || os.Getenv("M2M-ENVIRONMENT") == productionEnv {
		log.SetOutput(&lumberjack.Logger{
			Filename:   cfg.Config.Logging.File,
			MaxSize:    1,
			MaxBackups: 14,
			MaxAge:     28,
		})
	} else {
		log.SetOutput(os.Stdout)
	}

	log.Println("Configurações carregadas!")
}

func main() {
	app.Bootstrap()
}
