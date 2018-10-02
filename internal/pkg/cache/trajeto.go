package cache

import (
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"gopkg.in/mgo.v2/bson"
)

//CacheTrajeto -
var cacheTrajeto *Trajeto

//Trajeto -
type Trajeto struct {
	iniciado          bool
	trajetoRepository *repository.TrajetoRepository
	cache             map[string]model.Trajeto
	lock              chan interface{}
}

func newTrajeto(trajetoRepository *repository.TrajetoRepository) (*Trajeto, error) {
	m := new(Trajeto)
	lockSize := 1
	m.lock = make(chan interface{}, lockSize*2)
	for i := 0; i < lockSize; i++ {
		v := new(interface{})
		m.lock <- v
	}
	m.trajetoRepository = trajetoRepository
	err := m.criar()
	m.manterCacheAtualizado()
	return m, err
}

//GetTrajeto retorna instancia funcional de cache de trajeto
func GetTrajeto(trajetoRepository *repository.TrajetoRepository) (*Trajeto, error) {
	var err error
	if cacheTrajeto == nil {
		m, err := newTrajeto(trajetoRepository)
		cacheTrajeto = m
		if err != nil {
			return nil, err
		}
	}
	return cacheTrajeto, err
}

func (p *Trajeto) manterCacheAtualizado() {
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

func (p *Trajeto) atualizar() error {
	listaTrajetos, err := p.trajetoRepository.ListarTrajetos()
	if err != nil {
		logger.Errorf("Trajetos: %v\n", err)
		return err
	}

	l := <-p.lock
	defer p.lockRelease(l)

	for _, t := range listaTrajetos {
		p.cache[t.ID.Hex()] = t
	}

	logger.Debugf("Trajetos Atualizado: %v\n", len(p.cache))
	return err
}

func (p *Trajeto) criar() error {
	var err error

	l := <-p.lock
	defer p.lockRelease(l)

	cache, err := p.trajetoRepository.CarregarMapaTrajetos()
	if err == nil {
		p.cache = cache
		p.iniciado = true
	}
	return err
}

//Get -
func (p *Trajeto) Get(id string) (model.Trajeto, error) {
	if v, k := p.cache[id]; k {
		return v, nil
	}
	return p.find(id)
}

//find -
func (p *Trajeto) find(id string) (model.Trajeto, error) {
	l := <-p.lock
	defer p.lockRelease(l)
	if v, k := p.cache[id]; k {
		return v, nil
	}
	trajeto, err := p.trajetoRepository.ConsultarPorID(bson.ObjectIdHex(id))
	if trajeto.ID.Valid() {
		p.cache[id] = trajeto
	}
	return trajeto, err
}

func (p *Trajeto) lockRelease(l interface{}) {
	p.lock <- l
}
