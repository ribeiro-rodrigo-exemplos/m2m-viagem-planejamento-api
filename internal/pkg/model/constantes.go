package model

import (
	"strconv"
	"time"
)

//TiposDia -
var TiposDia = &TiposDiaConst{
	SEGUNDA:  TipoDiaConst{Dia: "2"},
	TERCA:    TipoDiaConst{Dia: "3"},
	QUARTA:   TipoDiaConst{Dia: "4"},
	QUINTA:   TipoDiaConst{Dia: "5"},
	SEXTA:    TipoDiaConst{Dia: "6"},
	SABADO:   TipoDiaConst{Dia: "S"},
	DOMINGO:  TipoDiaConst{Dia: "D"},
	ESPECIAL: TipoDiaConst{Dia: "O"},
	FERIADO:  TipoDiaConst{Dia: "F"},
	UTIL:     TipoDiaConst{Dia: "U"},
}

//TiposDiaConst -
type TiposDiaConst struct {
	SEGUNDA  TipoDiaConst
	TERCA    TipoDiaConst
	QUARTA   TipoDiaConst
	QUINTA   TipoDiaConst
	SEXTA    TipoDiaConst
	SABADO   TipoDiaConst
	DOMINGO  TipoDiaConst
	ESPECIAL TipoDiaConst
	FERIADO  TipoDiaConst
	UTIL     TipoDiaConst
}

//FromDate -
func (td *TiposDiaConst) FromDate(date *time.Time, others []string) []string {
	dias := []string{}
	if others != nil {
		dias = append(dias, others...)
	}
	weekday := (int(date.Weekday()) + 1)
	day := strconv.Itoa(weekday)

	if weekday != 1 && weekday != 7 {
		switch day {
		case "2":
			dias = append(dias, td.SEGUNDA.Dia)
		case "3":
			dias = append(dias, td.TERCA.Dia)
		case "4":
			dias = append(dias, td.QUARTA.Dia)
		case "5":
			dias = append(dias, td.QUINTA.Dia)
		case "6":
			dias = append(dias, td.SEXTA.Dia)
		}
		dias = append(dias, td.UTIL.Dia)
	} else if weekday == 1 {
		dias = append(dias, td.DOMINGO.Dia)
	} else if weekday == 7 {
		dias = append(dias, td.SABADO.Dia)
	}
	return dias
}

//TipoDiaConst -
type TipoDiaConst struct {
	Dia string
}

//ViagemEstado -
var ViagemEstado = ViagemEstadoConst{
	NovaViagem:          1,
	EmPreparacao:        2,
	ViagemAberta:        3,
	ViagemFechada:       4,
	ViagemCancelada:     5,
	ViagemPendente:      6,
	DeslocamentoEmCerca: 7,
}

//ViagemEstadoConst -
type ViagemEstadoConst struct {
	NovaViagem          int32
	EmPreparacao        int32
	ViagemAberta        int32
	ViagemFechada       int32
	ViagemCancelada     int32
	ViagemPendente      int32
	DeslocamentoEmCerca int32
}
