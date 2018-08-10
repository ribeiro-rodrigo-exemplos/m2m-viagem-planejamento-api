package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/jinzhu/configor"
)

// Config - configuracões da aplicação
var Config = Configuration{}
var loaded bool

// Configuration - type
type Configuration struct {
	RuntimeEnv string `default:"DEVELOPMENT"`
	ChaveJWT   string `required:"true"`

	Environment struct {
		Name string `required:"true"`
	}

	Server struct {
		Port string `default:"8080"`
	}

	HTTP struct {
		Request struct {
			MaxConcurrent int `default:"10"`
		}
		Transport struct {
			MaxIdleConnsPerHost int `default:"100"`
		}
	}

	Service struct {
		ViagemPlanejamento struct {
			MaxConcurrent int `default:"3"`
		}
	}

	MySQL struct {
		Host           string        `default:"localhost"`
		Port           int           `default:"3306"`
		User           string        `required:"true"`
		Password       string        `required:"true"`
		Database       string        `required:"true"`
		MaxIdleConns   int           `default:"3"`
		MaxOpenConns   int           `default:"10"`
		Reconnect      int           `default:"3"`
		ReconnectSleep time.Duration `default:"3"`
	}

	MongoDB struct {
		Host     string        `default:"localhost"`
		Port     int           `default:"27017"`
		Database string        `required:"true"`
		Timeout  time.Duration `default:"57"`
	}

	Hazelcast struct {
		Name string `required:"true"`
		Host string `default:"127.0.0.1"`
		Port string `default:"5701"`
	}

	Logging struct {
		File  string `required:"false"`
		Level map[string]string
	}
}

//InitConfig -
func InitConfig(configLocationFile string) {

	if !loaded {
		configLocation, environmentFlag := loadFlags()
		if configLocationFile != "" {
			configLocation = configLocationFile
		}

		Config = loadConfig(configLocation)
		Config.RuntimeEnv = environmentFlag
		loaded = true
	}

}

func loadFlags() (string, string) {
	//Para debugar com IDE definir caminho de configuração usando parâmetro -config-location=caminho_para_cfg.json
	configLocation := flag.String("config-location", "./configs/config.json", "a string")
	environment := flag.String("m2m-environment", "DEVELOPMENT", "a string")
	flag.Parse()

	return *configLocation, *environment
}

func loadConfig(configLocation string) Configuration {
	configApp := new(Configuration)

	err := configor.Load(configApp, configLocation)

	// fmt.Printf("%+v\n", configApp)
	if err != nil {
		fmt.Printf("Erro ao carregar configurações - %s\n", err)
	}
	return *configApp
}
