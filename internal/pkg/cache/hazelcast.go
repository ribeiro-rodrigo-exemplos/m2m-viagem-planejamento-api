package cache

import (
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
)

var client hazelcast.IHazelcastInstance

//InsereIntervalo - responsável por inserir no Hazelcast o novo valor de intervalo de métricas
func InsereIntervalo(clienteID int64, intervalo int) error {
	mp, err := client.GetMap(cfg.Config.Hazelcast.Name)

	if err != nil {
		_, err = mp.Put(clienteID, intervalo)
	}

	return err
}

//PegaIntervalo - responsável por pegar o valor de intervalo de métricas no Hazelcast
func PegaIntervalo(clienteID int64) (interface{}, error) {
	var intervalo interface{}
	mp, err := client.GetMap(cfg.Config.Hazelcast.Name)

	if err != nil {
		intervalo, err = mp.Get(clienteID)
	}

	return intervalo, err
}

func configuraClient() hazelcast.IHazelcastInstance {

	configuracao := hazelcast.NewHazelcastConfig()
	configuracao.ClientNetworkConfig().AddAddress(cfg.Config.Hazelcast.Host + ":" + cfg.Config.Hazelcast.Port)

	client, err := hazelcast.NewHazelcastClientWithConfig(configuracao)

	if err != nil {
		logger.Errorf("Erro ao configurar o client do Hazelcast %s\n", err)
	}

	return client
}
