package util

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

//ObterTimezoneTime - obtem time UTC
func ObterTimezoneTime(l *time.Location, dataHora string) (*time.Time, error) {

	timeParse, err := time.ParseInLocation("2006-01-02 15:04:05", dataHora, l)

	if err != nil {
		return &timeParse, err
	}

	return &timeParse, err
}

//ObterMes - obtem o mês do ano
func ObterMes(dataHora string) (string, error) {

	timeParse, err := time.Parse("2006-01-02", dataHora)

	if err != nil {
		return "", errors.New(fmt.Sprint("Erro ao formatar data da transmissao", err))
	}

	mes := fmt.Sprintf("%04d-%02d", timeParse.Year(), int(timeParse.Month()))

	return mes, nil
}

//ObterData - obtem a data
func ObterData(dataHora string) string {
	if strings.Contains(dataHora, "T") {
		return strings.Split(dataHora, "T")[0]
	}
	return strings.Split(dataHora, " ")[0]
}

//ObterDataAtual - obtem a data atual
func ObterDataAtual() string {
	dataHoraAtual := time.Now().Local()
	return dataHoraAtual.Format("2006-01-02")
}

//FormatarHM retornar hh:mm:ss
func FormatarHM(t time.Time) string {
	timeFormat := t.Format("15:04")
	return timeFormat
}

//FormatarHMS retornar hh:mm:ss
func FormatarHMS(t *time.Time) string {
	timeFormat := t.Format("15:04:05")
	return timeFormat
}

//FormatarAMDHMS retornar aaa-mm-dd hh:mm:ss
func FormatarAMDHMS(t *time.Time) string {
	timeFormat := t.Format("2006-01-02 15:04:05")
	return timeFormat
}

//DuracaoEFormatacao -
func DuracaoEFormatacao(inicio *time.Time, fim *time.Time) (duration time.Duration, formatacao string) {
	duration = fim.Sub(*inicio)
	durationRounded := duration.Round(time.Second)
	var neg bool
	if durationRounded < 0 {
		durationRounded = durationRounded * -1
		neg = true
	}
	h := durationRounded / time.Hour
	durationRounded -= h * time.Hour
	m := durationRounded / time.Minute
	durationRounded -= m * time.Minute
	s := durationRounded / time.Second
	if !neg {
		formatacao = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	} else {
		formatacao = fmt.Sprintf("-%02d:%02d:%02d", h, m, s)
	}

	return
}

//DuracaoEFormatacaoMinutos -
func DuracaoEFormatacaoMinutos(inicio *time.Time, fim *time.Time) (duration time.Duration, formatacao string) {
	duration = fim.Sub(*inicio)
	durationRounded := duration.Round(time.Minute)
	var neg bool
	if durationRounded < 0 {
		durationRounded = durationRounded * -1
		neg = true
	}
	h := durationRounded / time.Hour
	durationRounded -= h * time.Hour
	m := durationRounded / time.Minute
	if !neg {
		formatacao = fmt.Sprintf("%02d:%02d", h, m)
	} else {
		formatacao = fmt.Sprintf("-%02d:%02d", h, m)
	}

	return
}

//DuracaoEFormatacaoMinutosTrunc -
func DuracaoEFormatacaoMinutosTrunc(inicio time.Time, fim time.Time) (duration time.Duration, formatacao string) {
	duration = fim.Truncate(60 * time.Second).Sub(inicio.Truncate(60 * time.Second))
	durationRounded := duration.Round(time.Minute)
	var neg bool
	if durationRounded < 0 {
		durationRounded = durationRounded * -1
		neg = true
	}
	h := durationRounded / time.Hour
	durationRounded -= h * time.Hour
	m := durationRounded / time.Minute
	if !neg {
		formatacao = fmt.Sprintf("%02d:%02d", h, m)
	} else {
		formatacao = fmt.Sprintf("-%02d:%02d", h, m)
	}

	return
}

//SplitDiasPeriodo -
func SplitDiasPeriodo(periodoInicial Periodo) []Periodo {
	inicio := periodoInicial.Inicio
	fim := periodoInicial.Fim
	diff := fim.Sub(*inicio)
	//fmt.Printf("DURATION: %v", diff)
	if diff.Hours() < 24 && inicio.Day() == fim.Day() {
		periodos := []Periodo{periodoInicial}
		return periodos
	}

	dataAtual := inicio
	var novaData *time.Time
	periodos := make([]Periodo, 0, 10)

	for i := 0; true; i++ {
		novaData = ArredondarFimDia(dataAtual)
		periodo := &Periodo{
			Inicio: dataAtual,
			Fim:    novaData,
		}
		periodos = append(periodos, *periodo)
		dataAtualAux := novaData.Add(1 * time.Second)
		dataAtual = &dataAtualAux

		if fim.Sub(*dataAtual).Hours() < 24 && dataAtual.Day() == fim.Day() {
			periodo := Periodo{
				Inicio: dataAtual,
				Fim:    fim,
			}
			periodos = append(periodos, periodo)
			break
		}

	}
	return periodos
}

//ArredondarInicioDia -
func ArredondarInicioDia(t *time.Time) *time.Time {
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &rounded
}

//ArredondarFimDia -
func ArredondarFimDia(t *time.Time) *time.Time {
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
	return &rounded
}

//Concatenar -
func Concatenar(data time.Time, hora time.Time, loc *time.Location) (novaDataHora time.Time) {
	if loc != nil {
		novaDataHora = time.Date(data.Year(), data.Month(), data.Day(), hora.Hour(), hora.Minute(), hora.Second(), 0, loc)
	}
	return
}

//Periodo -
type Periodo struct {
	Inicio *time.Time
	Fim    *time.Time
}
