package util

import (
	"testing"
	"time"
)

func TestFormatarDataComTimezone(t *testing.T) {
	var err error

	tmzs := []struct {
		loc       string
		resultado string
	}{
		{"UTC", "2018-07-27 12:05:18"},
		{"America/Sao_Paulo", "2018-07-27 09:05:18"},
		{"America/Maceio", "2018-07-27 09:05:18"},
		{"America/Fortaleza", "2018-07-27 09:05:18"},
		{"America/Belem", "2018-07-27 09:05:18"},
		{"Europe/Lisbon", "2018-07-27 13:05:18"},
		{"America/Santarem", "2018-07-27 09:05:18"},
		{"America/Cuiaba", "2018-07-27 08:05:18"},
		{"America/Bahia", "2018-07-27 09:05:18"},
		{"America/Rio_Branco", "2018-07-27 07:05:18"},
		{"America/Recife", "2018-07-27 09:05:18"},
	}

	dt, err := time.ParseInLocation("2006-01-02 15:04:05", "2018-07-27 12:05:18", time.UTC)
	if err != nil {
		t.Error(err)
	}
	for _, tmz := range tmzs {
		loc, err := time.LoadLocation(tmz.loc)
		if err != nil {
			t.Error(err)
			continue
		}

		d, err := FormatarDataComTimezone(dt, loc)
		if err != nil {
			t.Error(err)
			continue
		}
		t.Logf("%v - %v\n", d, tmz.loc)

		if d != tmz.resultado {
			t.Errorf("Data esperada %q é diferente de %q\n", tmz.resultado, d)
		}

	}

}

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

	d, f := DuracaoEFormatacao(&inicio, &fim)

	if formatacaoEsperada != *f {
		t.Errorf("Formatação esperada %s é diferente de %v\n", formatacaoEsperada, f)
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

	d, f := DuracaoEFormatacao(&inicio, &fim)

	if formatacaoEsperada != *f {
		t.Errorf("Formatação esperada %s é diferente de %v\n", formatacaoEsperada, f)
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

	d, f := DuracaoEFormatacao(&inicio, &fim)

	if formatacaoEsperada != *f {
		t.Errorf("Formatação esperada %s é diferente de %v\n", formatacaoEsperada, f)
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

	d, f := DuracaoEFormatacao(&inicio, &fim)

	if formatacaoEsperada != *f {
		t.Errorf("Formatação esperada %s é diferente de %v\n", formatacaoEsperada, f)
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

	d, f := DuracaoEFormatacao(&inicio, &fim)

	if formatacaoEsperada != *f {
		t.Errorf("Formatação esperada %s é diferente de %v\n", formatacaoEsperada, f)
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

	d, f := DuracaoEFormatacao(&inicio, &fim)

	if formatacaoEsperada != *f {
		t.Errorf("Formatação esperada %s é diferente de %v\n", formatacaoEsperada, f)
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

	d, f := DuracaoEFormatacao(&inicio, &fim)

	if formatacaoEsperada != *f {
		t.Errorf("Formatação esperada %s é diferente de %v\n", formatacaoEsperada, f)
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

	d, f := DuracaoEFormatacao(&inicio, &fim)

	if formatacaoEsperada != *f {
		t.Errorf("Formatação esperada %s é diferente de %v\n", formatacaoEsperada, f)
	}

	if d.Seconds() != float64(duracaoEsperada) {
		t.Errorf("Duração esperada %vs é diferente de %v\n", duracaoEsperada, d.Seconds())
	}
}

func TestDuracaoDeHorario1h2m3s(t *testing.T) {
	var err error
	duracaoEsperada := (1 * time.Hour) + (2 * time.Minute) + (3 * time.Second)
	hora := "01:02:03"

	d, err := DuracaoDeHorario(hora)

	if err != nil {
		t.Errorf("%v\n", err)
	}

	if d != duracaoEsperada {
		t.Errorf("Duração esperada %v é diferente de %v\n", duracaoEsperada, d)
	}
}

func TestDuracaoDeHorario1h0m3s(t *testing.T) {
	var err error
	duracaoEsperada := (1 * time.Hour) + (0 * time.Minute) + (3 * time.Second)
	hora := "01:00:03"

	d, err := DuracaoDeHorario(hora)

	// fmt
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if d != duracaoEsperada {
		t.Errorf("Duração esperada %v é diferente de %v\n", duracaoEsperada, d)
	}
}

func TestDuracaoDeHorario3s(t *testing.T) {
	var err error
	duracaoEsperada := (0 * time.Hour) + (0 * time.Minute) + (3 * time.Second)
	hora := "00:00:03"

	d, err := DuracaoDeHorario(hora)

	// fmt
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if d != duracaoEsperada {
		t.Errorf("Duração esperada %v é diferente de %v\n", duracaoEsperada, d)
	}
}

func TestDuracaoDeHorario1h2m3s14ms(t *testing.T) {
	var err error
	duracaoEsperada := (1 * time.Hour) + (2 * time.Minute) + (3 * time.Second) + (14 * time.Millisecond)
	hora := "01:02:03:14"

	d, err := DuracaoDeHorario(hora)

	// fmt
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if d != duracaoEsperada {
		t.Errorf("Duração esperada %v é diferente de %v\n", duracaoEsperada, d)
	}
}

func TestDuracaoDeHorario14ms(t *testing.T) {
	var err error
	duracaoEsperada := (0 * time.Hour) + (0 * time.Minute) + (0 * time.Second) + (14 * time.Millisecond)
	hora := "00:00:00:14"

	d, err := DuracaoDeHorario(hora)

	// fmt
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if d != duracaoEsperada {
		t.Errorf("Duração esperada %v é diferente de %v\n", duracaoEsperada, d)
	}
}

func TestDuracaoDeHorarioNanosegundosDesconsiderados(t *testing.T) {
	var err error
	duracaoEsperada := (0 * time.Hour) + (0 * time.Minute) + (0 * time.Second) + (0 * time.Millisecond) + (0 * time.Nanosecond)
	hora := "00:00:00:00:14"

	d, err := DuracaoDeHorario(hora)

	// fmt
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if d != duracaoEsperada {
		t.Errorf("Duração esperada %v é diferente de %v\n", duracaoEsperada, d)
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
		Inicio: &inicio,
		Fim:    &fim,
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
		Inicio: &inicio,
		Fim:    &fim,
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
		Inicio: &inicio,
		Fim:    &fim,
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
		Inicio: &inicio,
		Fim:    &fim,
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
		Inicio: &inicio,
		Fim:    &fim,
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
		Inicio: &inicio,
		Fim:    &fim,
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
