package cache

import (
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"gopkg.in/mgo.v2/bson"
)

//CachePontoInteresse -
var cachePontoInteresse *PontoInteresse

//PontoInteresse -
type PontoInteresse struct {
	iniciado                 bool
	pontoInteresseRepository *repository.PontoInteresseRepository
	cache                    map[bson.ObjectId]*model.PontoInteresse
	lock                     chan interface{}
}

func newPontoInteresse(pontoInteresseRepository *repository.PontoInteresseRepository) (*PontoInteresse, error) {
	m := new(PontoInteresse)
	lockSize := 1
	m.lock = make(chan interface{}, lockSize*2)
	for i := 0; i < lockSize; i++ {
		v := new(interface{})
		m.lock <- v
	}
	m.pontoInteresseRepository = pontoInteresseRepository
	err := m.criar()
	m.manterCacheAtualizado()
	return m, err
}

//GetPontoInteresse retorna instancia funcional de cache de pontoInteresse
func GetPontoInteresse(pontoInteresseRepository *repository.PontoInteresseRepository) (*PontoInteresse, error) {
	var err error
	if cachePontoInteresse == nil {
		m, err := newPontoInteresse(pontoInteresseRepository)
		cachePontoInteresse = m
		if err != nil {
			return nil, err
		}
	}
	return cachePontoInteresse, err
}

func (p *PontoInteresse) manterCacheAtualizado() {
	go func() {
		if !p.iniciado {
			for i := 0; i < 3; i++ {
				time.Sleep(5 * time.Second)
				p.criar()
			}
			if !p.iniciado {
				return
			}
		}
		for {
			select {
			case <-time.After(60 * time.Second):
				p.atualizar()
			}
		}
	}()
}

func (p *PontoInteresse) atualizar() error {
	//Atualiza nome dos Pontos de Interesse
	err := p.criar()
	if err != nil {
		logger.Errorf("PontoInteresses: %v\n", err)
	} else {
		logger.Debugf("PontoInteresses Atualizado: %v\n", len(p.cache))
	}
	return err
}

func (p *PontoInteresse) criar() error {
	//Consultar IDs de Ponto Final de Trajetos de Linhas
	//Consultar Pontos de Intersse Ativos dos IDs recuperados
	//Montar cache
	cache, err := p.pontoInteresseRepository.CarregarMapaPontoInteresses(p.keys())
	if err == nil {
		p.cache = cache
		p.iniciado = true
	}
	return err
}

func (p *PontoInteresse) keys() []bson.ObjectId {
	keys := make([]bson.ObjectId, len(p.cache))
	i := 0
	for k := range p.cache {
		keys[i] = k
		i++
	}
	return keys
}

//Get -
func (p *PontoInteresse) Get(id bson.ObjectId) (*model.PontoInteresse, error) {
	if v, k := p.cache[id]; k {
		return v, nil
	}
	return p.find(id)
}

//find -
func (p *PontoInteresse) find(id bson.ObjectId) (*model.PontoInteresse, error) {
	l := <-p.lock
	defer p.lockRelease(l)
	if v, k := p.cache[id]; k {
		return v, nil
	}
	pontoInteresse, err := p.pontoInteresseRepository.ConsultarPorID(id)
	return pontoInteresse, err
}

func (p *PontoInteresse) lockRelease(l interface{}) {
	p.lock <- l
}
