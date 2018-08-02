package util

import (
	"testing"
	"time"
)

func TestDuracaoEFormatacaoSemDiferenca(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	duracaoEsperada := 0
	formatacaoEsperada := "00:00:00"

	d, f := DuracaoEFormatacao(inicio, fim)

	if formatacaoEsperada != f {
		t.Errorf("Formatação esperada %s é diferente de %q\n", formatacaoEsperada, f)
	}

	if d.Seconds() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Seconds())
	}
}

func TestDuracaoEFormatacao5Minutos(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:05:00")
	if err != nil {
		t.Error(err)
	}
	duracaoEsperada := 300
	formatacaoEsperada := "00:05:00"

	d, f := DuracaoEFormatacao(inicio, fim)

	if formatacaoEsperada != f {
		t.Errorf("Formatação esperada %s é diferente de %q\n", formatacaoEsperada, f)
	}

	if d.Seconds() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Seconds())
	}
}

func TestDuracaoEFormatacao02h00m15s(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 14:00:15")
	if err != nil {
		t.Error(err)
	}
	duracaoEsperada := 7215
	formatacaoEsperada := "02:00:15"

	d, f := DuracaoEFormatacao(inicio, fim)

	if formatacaoEsperada != f {
		t.Errorf("Formatação esperada %s é diferente de %q\n", formatacaoEsperada, f)
	}

	if d.Seconds() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Seconds())
	}
}

func TestDuracaoEArredondamentoMinutos02h00m15s(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 14:00:15")
	if err != nil {
		t.Error(err)
	}
	duracaoEsperada := 120
	formatacaoEsperada := "02:00:15"

	d, f := DuracaoEFormatacao(inicio, fim)

	if formatacaoEsperada != f {
		t.Errorf("Formatação esperada %s é diferente de %q\n", formatacaoEsperada, f)
	}

	if d.Round(time.Minute).Minutes() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Round(time.Minute).Minutes())
	}
}

func TestDuracaoEFormatacao3horas(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 15:00:00")
	if err != nil {
		t.Error(err)
	}
	duracaoEsperada := 10800
	formatacaoEsperada := "03:00:00"

	d, f := DuracaoEFormatacao(inicio, fim)

	if formatacaoEsperada != f {
		t.Errorf("Formatação esperada %s é diferente de %q\n", formatacaoEsperada, f)
	}

	if d.Seconds() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Seconds())
	}
}

func TestDuracaoEFormatacao72h04m28s(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-30 12:04:28")
	if err != nil {
		t.Error(err)
	}
	duracaoEsperada := 259468
	formatacaoEsperada := "72:04:28"

	d, f := DuracaoEFormatacao(inicio, fim)

	if formatacaoEsperada != f {
		t.Errorf("Formatação esperada %s é diferente de %q\n", formatacaoEsperada, f)
	}

	if d.Seconds() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Seconds())
	}
}

func TestDuracaoEFormatacaoNegativo00h05m27s(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 19:05:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 18:59:33")
	if err != nil {
		t.Error(err)
	}
	duracaoEsperada := -327
	formatacaoEsperada := "-00:05:27"

	d, f := DuracaoEFormatacao(inicio, fim)

	if formatacaoEsperada != f {
		t.Errorf("Formatação esperada %s é diferente de %q\n", formatacaoEsperada, f)
	}

	if d.Seconds() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Seconds())
	}
}

func TestDuracaoEFormatacaoNegativo24h00m10s(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 19:05:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-26 19:04:50")
	if err != nil {
		t.Error(err)
	}
	duracaoEsperada := -86410
	formatacaoEsperada := "-24:00:10"

	d, f := DuracaoEFormatacao(inicio, fim)

	if formatacaoEsperada != f {
		t.Errorf("Formatação esperada %s é diferente de %q\n", formatacaoEsperada, f)
	}

	if d.Seconds() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Seconds())
	}
}

func TestSplitPeriodoPorHoraMesmoDia(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	periodo := Periodo{
		Inicio: inicio,
		Fim:    fim,
	}

	periodos := SplitDiasPeriodo(periodo)

	if periodos == nil {
		t.Errorf("Periodos não pode ser nulo\n")
		return
	}
	if len(periodos) < 1 {
		t.Errorf("Periodos %+v não pode ser vazio.\n", periodos)
	}
	if len(periodos) != 1 {
		t.Errorf("Periodos %+v deve ter 1 elemento\n", periodos)
	}
	for _, p := range periodos {
		t.Logf("%+v\n", p)
	}

}

func TestSplitPeriodoPorHoraDiaSeguinte(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-28 12:00:00")
	if err != nil {
		t.Error(err)
	}
	periodo := Periodo{
		Inicio: inicio,
		Fim:    fim,
	}

	periodos := SplitDiasPeriodo(periodo)

	if periodos == nil {
		t.Errorf("Periodos não pode ser nulo\n")
		return
	}
	if len(periodos) < 1 {
		t.Errorf("Periodos %v não pode ser vazio.\n", periodos)
	}
	if len(periodos) != 2 {
		t.Errorf("Periodos %v deve ter 2 elementos\n", periodos)
	}
	for _, p := range periodos {
		t.Logf("%+v\n", p)
	}

}

func TestSplitPeriodoPorHoraDiaSeguinteMenosDe24Horas(t *testing.T) {
	t.Log("TestSplitPeriodoPorHoraDiaSeguinteMenosDe24Horas")

	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 18:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-28 17:59:59")
	if err != nil {
		t.Error(err)
	}
	periodo := Periodo{
		Inicio: inicio,
		Fim:    fim,
	}

	periodos := SplitDiasPeriodo(periodo)

	if periodos == nil {
		t.Errorf("Periodos não pode ser nulo\n")
		return
	}
	if len(periodos) < 1 {
		t.Errorf("Periodos %+v não pode ser vazio.\n", periodos)
	}
	if len(periodos) != 2 {
		t.Errorf("Periodos %+v deve ter 2 elementos e não %d\n", periodos, len(periodos))
	}
	for _, p := range periodos {
		t.Logf("%+v\n", p)
	}

}

func TestSplitPeriodoPorHora30Dias(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-27 12:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-08-26 12:00:00")
	if err != nil {
		t.Error(err)
	}
	periodo := Periodo{
		Inicio: inicio,
		Fim:    fim,
	}

	periodos := SplitDiasPeriodo(periodo)

	if periodos == nil {
		t.Errorf("Periodos não pode ser nulo\n")
		return
	}

	if len(periodos) < 1 {
		t.Errorf("Periodos %v não pode ser vazio.\n", periodos)
	}
	if len(periodos) != 31 {
		t.Errorf("Periodos \n%v\n deve ter 31 elementos e não %d\n", periodos, len(periodos))
	}
	for _, p := range periodos {
		t.Logf("%+v\n", p)
	}

}

func TestSplitPeriodoPorHora28Dias(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-07-01 18:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-28 20:00:00")
	if err != nil {
		t.Error(err)
	}
	periodo := Periodo{
		Inicio: inicio,
		Fim:    fim,
	}

	periodos := SplitDiasPeriodo(periodo)

	if periodos == nil {
		t.Errorf("Periodos não pode ser nulo\n")
		return
	}

	if len(periodos) < 1 {
		t.Errorf("Periodos %v não pode ser vazio.\n", periodos)
	}
	if len(periodos) != 28 {
		t.Errorf("Periodos \n%v\n deve ter 28 elementos e não %d\n", periodos, len(periodos))
	}
	for _, p := range periodos {
		t.Logf("%+v\n", p)
	}

}

func TestSplitPeriodoPorHora31DiasMesSeguinte(t *testing.T) {
	var err error
	inicio, err := time.Parse("2006-01-02 15:04:05", "2018-06-15 18:00:00")
	if err != nil {
		t.Error(err)
	}
	fim, err := time.Parse("2006-01-02 15:04:05", "2018-07-15 20:00:00")
	if err != nil {
		t.Error(err)
	}
	periodo := Periodo{
		Inicio: inicio,
		Fim:    fim,
	}

	periodos := SplitDiasPeriodo(periodo)

	if periodos == nil {
		t.Errorf("Periodos não pode ser nulo\n")
		return
	}

	if len(periodos) < 1 {
		t.Errorf("Periodos %v não pode ser vazio.\n", periodos)
	}
	if len(periodos) != 31 {
		t.Errorf("Periodos \n%v\n deve ter 31 elementos e não %d\n", periodos, len(periodos))
	}
	for _, p := range periodos {
		t.Logf("%+v\n", p)
	}

}
