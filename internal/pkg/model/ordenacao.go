package model

import (
	"sort"
)

//OrdenarViagemExecutadaPorData ordena viagem executada por data
func OrdenarViagemExecutadaPorData(lista []*ViagemExecutada) {
	sort.SliceStable(lista, func(i, j int) bool {
		vgexi := lista[i]
		vgexj := lista[j]
		return vgexi.Executada.DataInicio.Before(*vgexj.Executada.DataInicio)
	})
}
