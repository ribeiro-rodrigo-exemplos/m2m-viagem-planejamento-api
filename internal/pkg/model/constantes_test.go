package model

import (
	"fmt"
	"testing"
	"time"
)

func TestTipoDiaNaoNulo(t *testing.T) {

	tiposDeDia := TiposDia.FromDate(time.Now(), []string{"O", "F"})

	fmt.Println(tiposDeDia)

	if tiposDeDia == nil {
		t.Errorf("TiposDeDia não pode ser nula\n")
		return
	}

	if len(tiposDeDia) < 1 {
		t.Errorf("TiposDeDia não pode ser vazia\n")
	}
}
