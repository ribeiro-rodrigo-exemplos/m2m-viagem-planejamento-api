//Package dto precisa ser refatorada para extrair comportamentos e permitir composições ad-hoc de ordenações
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

//OrdenarViagemPorTabela ordena viagem por tabela
func OrdenarViagemPorTabela(lista []*ViagemDTO) {
	sort.SliceStable(lista, func(i, j int) bool {
		vgi := lista[i]
		vgj := lista[j]

		if vgi.Trajeto.Linha.Numero < vgj.Trajeto.Linha.Numero {
			return true
		} else if vgi.Trajeto.Linha.Numero > vgj.Trajeto.Linha.Numero {
			return false
		}

		var tabelaI string
		var tabelaJ string

		if vgi.NmTabela != nil {
			tabelaI = *vgi.NmTabela
		}
		if vgj.NmTabela != nil {
			tabelaJ = *vgj.NmTabela
		}

		if tabelaI < tabelaJ {
			return true
		} else if tabelaI > tabelaJ {
			return false
		}

		return vgi.PartidaOrdenacao.Before(*vgj.PartidaOrdenacao)
	})
}
