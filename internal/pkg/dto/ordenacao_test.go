package dto

import (
	"testing"
	"time"
)

//OrdenarViagemPorData ordena viagem executada por data
func TestOrdenarViagemPorData(t *testing.T) {
	var lista []*ViagemDTO

	data1 := time.Now()
	data2 := data1.Add(-1 * time.Minute)
	data3 := data1.Add(1 * time.Minute)

	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data1, Status: &StatusViagem.NaoRealizada})
	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data2, Status: &StatusViagem.RealizadaNaoPlanejada})
	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data3, Status: &StatusViagem.Cancelada})

	vg1 := lista[0]
	vg2 := lista[1]
	vg3 := lista[2]

	if vg1.Status != &StatusViagem.NaoRealizada || vg2.Status != &StatusViagem.RealizadaNaoPlanejada || vg3.Status != &StatusViagem.Cancelada {
		t.Errorf("Dados entrada. Viagens em ordem errada. Status \n")
	}
	if vg1.PartidaOrdenacao != &data1 || vg2.PartidaOrdenacao != &data2 || vg3.PartidaOrdenacao != &data3 {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao\n")
	}
	if vg1.PartidaOrdenacao.Before(*vg2.PartidaOrdenacao) {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao: Esperado %v > %v\n", vg1.PartidaOrdenacao.Format("2006-01-02 15:04"), vg2.PartidaOrdenacao.Format("2006-01-02 15:04"))
	}
	if vg3.PartidaOrdenacao.Before(*vg1.PartidaOrdenacao) {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao: Esperado %v > %v\n", vg3.PartidaOrdenacao.Format("2006-01-02 15:04"), vg1.PartidaOrdenacao.Format("2006-01-02 15:04"))
	}

	t.Logf("ANTES : ")
	for i, v := range lista {
		t.Logf("\t[%d] %v {%v}\n", i, v.PartidaOrdenacao.Format("2006-01-02 15:04"), v.Status)
	}

	OrdenarViagemPorData(lista)

	vg1 = lista[0]
	vg2 = lista[1]
	vg3 = lista[2]

	t.Logf("DEPOIS : ")
	for i, v := range lista {
		t.Logf("\t[%d] %v {%v}\n", i, v.PartidaOrdenacao.Format("2006-01-02 15:04"), v.Status)
	}

	if vg1.Status != &StatusViagem.RealizadaNaoPlanejada || vg2.Status != &StatusViagem.NaoRealizada || vg3.Status != &StatusViagem.Cancelada {
		t.Errorf("Validação. Viagens em ordem errada. Status \n")
	}
	if vg1.PartidaOrdenacao != &data2 || vg2.PartidaOrdenacao != &data1 || vg3.PartidaOrdenacao != &data3 {
		t.Errorf("Validação. Viagens em ordem errada. PartidaOrdenacao\n")
	}
	if vg2.PartidaOrdenacao.Before(*vg1.PartidaOrdenacao) {
		t.Errorf("Validação. Viagens em ordem errada. PartidaOrdenacao: Esperado %v > %v\n", vg2.PartidaOrdenacao.Format("2006-01-02 15:04"), vg1.PartidaOrdenacao.Format("2006-01-02 15:04"))
	}
	if vg3.PartidaOrdenacao.Before(*vg2.PartidaOrdenacao) {
		t.Errorf("Validação. Viagens em ordem errada. PartidaOrdenacao: Esperado %v > %v\n", vg3.PartidaOrdenacao.Format("2006-01-02 15:04"), vg2.PartidaOrdenacao.Format("2006-01-02 15:04"))
	}

}

//TestOrdenarViagemPorLinha ordena viagem executada por data
func TestOrdenarViagemPorLinha(t *testing.T) {
	var lista []*ViagemDTO

	data1 := time.Now()
	data2 := data1.Add(-1 * time.Minute)
	data3 := data1.Add(1 * time.Minute)

	trajeto1 := TrajetoDTO{Descricao: "t2", Linha: LinhaDTO{Numero: "A"}}
	trajeto2 := TrajetoDTO{Descricao: "t1", Linha: LinhaDTO{Numero: "B"}}
	trajeto3 := TrajetoDTO{Descricao: "t1", Linha: LinhaDTO{Numero: "A"}}

	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data1, Status: &StatusViagem.NaoRealizada, Trajeto: trajeto1})
	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data2, Status: &StatusViagem.RealizadaNaoPlanejada, Trajeto: trajeto2})
	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data3, Status: &StatusViagem.Cancelada, Trajeto: trajeto3})

	vg1 := lista[0]
	vg2 := lista[1]
	vg3 := lista[2]

	if vg1.Status != &StatusViagem.NaoRealizada || vg2.Status != &StatusViagem.RealizadaNaoPlanejada || vg3.Status != &StatusViagem.Cancelada {
		t.Errorf("Dados entrada. Viagens em ordem errada. Status \n")
	}
	if vg1.PartidaOrdenacao != &data1 || vg2.PartidaOrdenacao != &data2 || vg3.PartidaOrdenacao != &data3 {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao\n")
	}
	if vg1.PartidaOrdenacao.Before(*vg2.PartidaOrdenacao) {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao: Esperado %v > %v\n", vg1.PartidaOrdenacao.Format("2006-01-02 15:04"), vg2.PartidaOrdenacao.Format("2006-01-02 15:04"))
	}
	if vg3.PartidaOrdenacao.Before(*vg1.PartidaOrdenacao) {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao: Esperado %v > %v\n", vg3.PartidaOrdenacao.Format("2006-01-02 15:04"), vg1.PartidaOrdenacao.Format("2006-01-02 15:04"))
	}

	t.Logf("ANTES : ")
	for i, v := range lista {
		t.Logf("\t[%d] %v {%v} %v \n", i, v.PartidaOrdenacao.Format("2006-01-02 15:04"), v.Status, v.Trajeto)
	}

	OrdenarViagemPorLinha(lista)

	vg1 = lista[0]
	vg2 = lista[1]
	vg3 = lista[2]

	t.Logf("DEPOIS : ")
	for i, v := range lista {
		t.Logf("\t[%d] %v {%v} %v \n", i, v.PartidaOrdenacao.Format("2006-01-02 15:04"), v.Status, v.Trajeto)
	}

	if vg1.Status != &StatusViagem.Cancelada || vg2.Status != &StatusViagem.NaoRealizada || vg3.Status != &StatusViagem.RealizadaNaoPlanejada {
		t.Errorf("Validação. Viagens em ordem errada. Status \n")
	}

}

//TestOrdenarViagemPorTabela ordena viagem por tabela
func TestOrdenarViagemPorTabela(t *testing.T) {
	var lista []*ViagemDTO

	data1 := time.Now()
	data2 := data1.Add(-1 * time.Minute)
	data3 := data1.Add(1 * time.Minute)

	trajeto1 := TrajetoDTO{Descricao: "trajetoA", Linha: LinhaDTO{Numero: "A"}}
	trajeto2 := TrajetoDTO{Descricao: "trajetoA", Linha: LinhaDTO{Numero: "B"}}
	trajeto3 := TrajetoDTO{Descricao: "trajetoA", Linha: LinhaDTO{Numero: "A"}}

	tabela1 := "Tabela3"
	tabela2 := "Tabela2"
	tabela3 := "Tabela2"

	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data1, Status: &StatusViagem.NaoRealizada, Trajeto: trajeto1, NmTabela: &tabela1})
	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data2, Status: &StatusViagem.RealizadaNaoPlanejada, Trajeto: trajeto2, NmTabela: &tabela2})
	lista = append(lista, &ViagemDTO{PartidaOrdenacao: &data3, Status: &StatusViagem.Cancelada, Trajeto: trajeto3, NmTabela: &tabela3})

	vg1 := lista[0]
	vg2 := lista[1]
	vg3 := lista[2]

	if vg1.Status != &StatusViagem.NaoRealizada || vg2.Status != &StatusViagem.RealizadaNaoPlanejada || vg3.Status != &StatusViagem.Cancelada {
		t.Errorf("Dados entrada. Viagens em ordem errada. Status \n")
	}
	if vg1.PartidaOrdenacao != &data1 || vg2.PartidaOrdenacao != &data2 || vg3.PartidaOrdenacao != &data3 {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao\n")
	}
	if vg1.PartidaOrdenacao.Before(*vg2.PartidaOrdenacao) {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao: Esperado %v > %v\n", vg1.PartidaOrdenacao.Format("2006-01-02 15:04"), vg2.PartidaOrdenacao.Format("2006-01-02 15:04"))
	}
	if vg3.PartidaOrdenacao.Before(*vg1.PartidaOrdenacao) {
		t.Errorf("Dados entrada. Viagens em ordem errada. PartidaOrdenacao: Esperado %v > %v\n", vg3.PartidaOrdenacao.Format("2006-01-02 15:04"), vg1.PartidaOrdenacao.Format("2006-01-02 15:04"))
	}

	t.Logf("ANTES : \n")
	for i, v := range lista {
		t.Logf("\t[%d] %v {%v} %v %v \n", i, v.PartidaOrdenacao.Format("2006-01-02 15:04"), *v.Status, v.Trajeto.Linha.Numero, *v.NmTabela)
	}

	OrdenarViagemPorTabela(lista)

	vg1 = lista[0]
	vg2 = lista[1]
	vg3 = lista[2]

	t.Logf("DEPOIS : \n")
	for i, v := range lista {
		t.Logf("\t[%d] %v {%v} %v %v \n", i, v.PartidaOrdenacao.Format("2006-01-02 15:04"), *v.Status, v.Trajeto.Linha.Numero, *v.NmTabela)
	}

	if vg1.Status != &StatusViagem.Cancelada || vg2.Status != &StatusViagem.NaoRealizada || vg3.Status != &StatusViagem.RealizadaNaoPlanejada {
		t.Errorf("Validação. Viagens em ordem errada. Status \n")
	}

}
