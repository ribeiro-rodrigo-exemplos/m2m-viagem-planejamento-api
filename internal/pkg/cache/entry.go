package cache

import (
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
)

//AgrupamentoEntry -
type AgrupamentoEntry struct {
	Agrupamento *model.Agrupamento
	Linhas      []*model.Linha
	Trajetos    []*model.Trajeto
	TrajetosDTO []dto.TrajetoDTO
}

//NewAgrupamentoEntry -
func NewAgrupamentoEntry(a *model.Agrupamento) (agrupamentoEntry *AgrupamentoEntry) {
	agrupamentoEntry = new(AgrupamentoEntry)
	agrupamentoEntry.Agrupamento = a
	agrupamentoEntry.Linhas = []*model.Linha{}
	agrupamentoEntry.Trajetos = []*model.Trajeto{}
	agrupamentoEntry.TrajetosDTO = []dto.TrajetoDTO{}
	return
}

//AddLinha -
func (a *AgrupamentoEntry) AddLinha(l model.Linha) {
	a.Linhas = append(a.Linhas, &l)
}

//GetLinhas -
func (a *AgrupamentoEntry) GetLinhas() (linhas []*model.Linha) {
	linhas = a.Linhas
	return
}

//AddTrajeto -
func (a *AgrupamentoEntry) AddTrajeto(l *model.Trajeto) {
	a.Trajetos = append(a.Trajetos, l)
	a.TrajetosDTO = append(a.TrajetosDTO, dto.TrajetoDTO{
		ID:        &l.ID,
		Descricao: l.Nome,
		Linha: dto.LinhaDTO{
			Numero: l.Linha.Numero,
		},
	})

}

//GetTrajetos -
func (a *AgrupamentoEntry) GetTrajetos() (trajetos []*model.Trajeto) {
	trajetos = a.Trajetos
	return
}
