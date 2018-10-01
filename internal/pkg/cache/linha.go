package cache

import (
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"gopkg.in/mgo.v2/bson"
)

//CacheLinha -
var cacheLinha *Linha

//Linha -
type Linha struct {
	iniciado        bool
	linhaRepository *repository.LinhaRepository
	cache           map[bson.ObjectId]*model.Linha
}

func newLinha(linhaRepository *repository.LinhaRepository) (*Linha, error) {
	m := new(Linha)
	m.linhaRepository = linhaRepository
	err := m.criar()
	m.manterCacheAtualizado()
	return m, err
}

//GetLinha retorna instancia funcional de cache de linha
func GetLinha(linhaRepository *repository.LinhaRepository) (*Linha, error) {
	var err error
	if cacheLinha == nil {
		m, err := newLinha(linhaRepository)
		cacheLinha = m
		if err != nil {
			return nil, err
		}
	}
	return cacheLinha, err
}

func (m *Linha) manterCacheAtualizado() {
	go func() {
		if !m.iniciado {
			for i := 0; i < 3; i++ {
				time.Sleep(5 * time.Second)
				m.criar()
			}
			if !m.iniciado {
				return
			}
		}
		for {
			select {
			case <-time.After(60 * time.Second):
				m.atualizar()
			}
		}
	}()
}

func (m *Linha) atualizar() error {
	err := m.criar()
	if err != nil {
		logger.Errorf("Linhas: %v\n", err)
	} else {
		logger.Debugf("Linhas Atualizado: %v\n", len(m.cache))
	}
	return err
}

func (m *Linha) criar() (err error) {
	linhas, err := m.linhaRepository.Listar()
	if err == nil && linhas != nil {
		var mapaLinhas = make(map[bson.ObjectId]*model.Linha)
		for _, l := range linhas {
			mapaLinhas[l.ID] = &l
		}
		m.cache = mapaLinhas
		m.iniciado = true
	}
	return err
}

//Get -
func (m *Linha) Get(id bson.ObjectId) (linha *model.Linha, err error) {
	if v, k := m.cache[id]; k {
		value := *v
		linha = &value
	}
	return
}

//ListAll -
func (m *Linha) ListAll() (linhas []model.Linha, err error) {
	linhas = make([]model.Linha, len(m.cache))
	for _, v := range m.cache {
		linhas = append(linhas, *v)
	}
	return
}
