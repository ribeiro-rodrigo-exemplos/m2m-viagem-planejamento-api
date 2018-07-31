package cache

import (
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
)

//CacheCliente -
var cacheCliente *Cliente

//Cliente -
type Cliente struct {
	Cache map[int16]*model.Cliente
}

//GetCliente retorna instancia funcional de cache de cliente
func GetCliente(clienteRepository *repository.ClienteRepository) (*Cliente, error) {
	var err error
	if cacheCliente == nil {
		c := new(Cliente)
		c.Cache, err = clienteRepository.CarregarMapaClientes()
		if err != nil {
			return nil, err
		}
		cacheCliente = c
	}
	return cacheCliente, err
}
