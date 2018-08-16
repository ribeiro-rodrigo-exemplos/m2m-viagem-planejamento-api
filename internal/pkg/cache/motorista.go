package cache

import (
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
)

//CacheMotorista -
var cacheMotorista *Motorista

//Motorista -
type Motorista struct {
	iniciado            bool
	motoristaRepository *repository.MotoristaRepository
	Cache               map[string]*model.Motorista
}

func newMotorista(motoristaRepository *repository.MotoristaRepository) (*Motorista, error) {
	m := new(Motorista)
	m.motoristaRepository = motoristaRepository
	err := m.criar()
	m.manterCacheAtualizado()
	return m, err
}

//GetMotorista retorna instancia funcional de cache de motorista
func GetMotorista(motoristaRepository *repository.MotoristaRepository) (*Motorista, error) {
	var err error
	if cacheMotorista == nil {
		m, err := newMotorista(motoristaRepository)
		cacheMotorista = m
		if err != nil {
			return nil, err
		}
	}
	return cacheMotorista, err
}

func (m *Motorista) manterCacheAtualizado() {
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

func (m *Motorista) atualizar() error {
	err := m.criar()
	if err != nil {
		logger.Errorf("Motoristas: %v\n", err)
	} else {
		logger.Debugf("Motoristas Atualizado: %v\n", len(m.Cache))
	}
	return err
}

func (m *Motorista) criar() error {
	cache, err := m.motoristaRepository.CarregarMapaMotoristas()
	if err == nil {
		m.Cache = cache
		m.iniciado = true
	}
	return err
}
