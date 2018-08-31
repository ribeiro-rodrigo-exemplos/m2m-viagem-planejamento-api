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

//OrdenarViagemPorLinha ordena viagem por linha,trajeto, data
func OrdenarViagemPorLinha(lista []*ViagemDTO) {
	sort.SliceStable(lista, func(i, j int) bool {
		vgi := lista[i]
		vgj := lista[j]

		//string
		if vgi.Trajeto.Linha.Numero < vgj.Trajeto.Linha.Numero {
			return true
		} else if vgi.Trajeto.Linha.Numero > vgj.Trajeto.Linha.Numero {
			return false
		}
		if vgi.Trajeto.Descricao < vgj.Trajeto.Descricao {
			return true
		} else if vgi.Trajeto.Descricao > vgj.Trajeto.Descricao {
			return false
		}

		return vgi.PartidaOrdenacao.Before(*vgj.PartidaOrdenacao)
	})
}
