package repository

import (
	"fmt"
	"testing"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
)

func TestListar(t *testing.T) {
	cfg.InitConfig("../../../configs/config.json")
	logger = logging.NewLogger("repository", cfg.Config.Logging.Level)
	database.InitConfig()
	fmt.Println("TestListar")

	session, err := database.GetMongoDB()
	if err != nil {
		t.Errorf("Conexão banco de dados - %s\n", err)
	}

	c := NewLinhaRepository(session)
	linhas, err := c.Listar()
	if err != nil {
		t.Errorf("Erro ao Listar - %s\n", err)
	}
	if linhas == nil {
		t.Errorf("Lista de Linhas não pode ser nula\n")
	}
	if len(linhas) < 1 {
		t.Errorf("Lista de Linhas %v não pode ser vazia\n", linhas)
	}
	// for _, l := range linhas {
	// 	t.Logf("%+v", l)
	// 	fmt.Printf("%v - %v\n", l.ID, l.Nome)
	// }

	t.Logf("Linhas  %d\n", len(linhas))
}
