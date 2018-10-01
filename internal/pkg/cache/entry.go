package cache

import (
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
)

//AgrupamentoEntry -
type AgrupamentoEntry struct {
	Agrupamento *model.Agrupamento
	Linhas      []*model.Linha
	Trajetos    []*model.Trajeto
}

//NewAgrupamentoEntry -
func NewAgrupamentoEntry(a *model.Agrupamento) (agrupamentoEntry *AgrupamentoEntry) {
	agrupamentoEntry = new(AgrupamentoEntry)
	agrupamentoEntry.Agrupamento = a
	agrupamentoEntry.Linhas = []*model.Linha{}
	agrupamentoEntry.Trajetos = []*model.Trajeto{}
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
}

//GetTrajetos -
func (a *AgrupamentoEntry) GetTrajetos() (trajetos []*model.Trajeto) {
	trajetos = a.Trajetos
	return
}
