package dto

import (
	"sort"
)

//OrdenarViagemExecutadaPorData ordena viagem executada por data
func OrdenarViagemExecutadaPorData(lista []*ViagemDTO) {
	sort.SliceStable(lista, func(i, j int) bool {
		vgi := lista[i]
		vgj := lista[j]
		return vgi.PartidaOrdenacao.Before(vgj.PartidaOrdenacao)
	})
}
