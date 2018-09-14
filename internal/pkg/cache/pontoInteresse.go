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
	mapaPontoInteresses, err := p.pontoInteresseRepository.CarregarMapaPontoInteresses(p.keys())
	if err != nil {
		logger.Errorf("PontoInteresses: %v\n", err)
		return err
	}

	l := <-p.lock
	defer p.lockRelease(l)

	for k, v := range mapaPontoInteresses {
		p.cache[k] = v
	}

	logger.Debugf("PontoInteresses Atualizado: %v\n", len(p.cache))
	return err
}

func (p *PontoInteresse) criar() error {
	var err error

	l := <-p.lock
	defer p.lockRelease(l)

	listaIDsPontosFinal, err := p.pontoInteresseRepository.ListarIdentificacaoPontosFinal()
	cache, err := p.pontoInteresseRepository.CarregarMapaPontoInteresses(listaIDsPontosFinal)
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
		logger.Tracef("Valor recuperado em memória %v\n", v)
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
	if pontoInteresse != nil {
		p.cache[id] = pontoInteresse
	}
	return pontoInteresse, err
}

func (p *PontoInteresse) lockRelease(l interface{}) {
	p.lock <- l
}
