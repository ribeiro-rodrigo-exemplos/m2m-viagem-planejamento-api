package cache

import (
	"time"
)

//CacheAgrupamento -
var cacheAgrupamento *Agrupamento

//Agrupamento -
type Agrupamento struct {
	iniciado   bool
	linhaCache *Linha
	cache      map[int32]*AgrupamentoEntry
}

func newAgrupamento(linhaCache *Linha) (*Agrupamento, error) {
	m := new(Agrupamento)
	m.linhaCache = linhaCache
	err := m.criar()
	m.manterCacheAtualizado()
	return m, err
}

//GetAgrupamento retorna instancia funcional de cache de agrupamento
func GetAgrupamento(linhaCache *Linha) (*Agrupamento, error) {
	var err error
	if cacheAgrupamento == nil {
		m, err := newAgrupamento(linhaCache)
		cacheAgrupamento = m
		if err != nil {
			return nil, err
		}
	}
	return cacheAgrupamento, err
}

func (m *Agrupamento) manterCacheAtualizado() {
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

func (m *Agrupamento) atualizar() error {
	err := m.criar()
	if err != nil {
		logger.Errorf("Agrupamentos: %v\n", err)
	} else {
		logger.Debugf("Agrupamentos Atualizado: %v\n", len(m.cache))
	}
	return err
}

func (m *Agrupamento) criar() (err error) {
	linhas, err := m.linhaCache.ListAll()
	if err == nil && linhas != nil {
		var mapaAgrupamentos = make(map[int32]*AgrupamentoEntry)
		for _, l := range linhas {
			a := l.Agrupamento
			if a.ID < 1 {
				continue
			}

			entry, key := mapaAgrupamentos[a.ID]
			if !key {
				entry = NewAgrupamentoEntry(&a)
				mapaAgrupamentos[a.ID] = entry
			}

			entry.AddLinha(l)
			for _, t := range l.Trajetos {
				t.Linha = l
				entry.AddTrajeto(&t)
			}
		}
		m.cache = mapaAgrupamentos
		m.iniciado = true
	}
	return err
}

//Get -
func (m *Agrupamento) Get(id int32) (agrupamento *AgrupamentoEntry, err error) {
	if v, k := m.cache[id]; k {
		agrupamento = v
	}
	return
}
