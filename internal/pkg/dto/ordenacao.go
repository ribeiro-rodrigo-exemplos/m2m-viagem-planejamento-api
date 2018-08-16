package dto

import (
	"sort"
)

//OrdenarViagemPorData ordena viagem executada por data
func OrdenarViagemPorData(lista []*ViagemDTO) {
	sort.SliceStable(lista, func(i, j int) bool {
		vgi := lista[i]
		vgj := lista[j]
		return vgi.PartidaOrdenacao.Before(*vgj.PartidaOrdenacao)
	})
}
