package repository

import (
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//ViagemExecutadaRepository -
type ViagemExecutadaRepository struct {
	mongoDB *database.MongoDB
}

//NewViagemExecutadaRepository -
func NewViagemExecutadaRepository(mongoDB *database.MongoDB) *ViagemExecutadaRepository {
	c := new(ViagemExecutadaRepository)
	c.mongoDB = mongoDB
	return c
}

//ListarViagensPor filtro
func (v *ViagemExecutadaRepository) ListarViagensPor(filtro dto.FilterDTO) ([]*model.ViagemExecutada, error) {
	var err error

	session, err := v.mongoDB.GetSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	listarViagens := []*model.ViagemExecutada{}

	projecao := bson.M{
		"alocacao":                 1,
		"_id":                      1,
		"situacaoAtual":            1,
		"executada":                1,
		"porcentagemConclusao":     1,
		"planejada":                1,
		"tipoViagem":               1,
		"qntPassageiros":           1,
		"mensagemObs":              1,
		"partida.trajetoExecutado": 1,
		"codigoMotorista":          1,

		// "transmissoesRecebidas.transmissao.dataTransmissao": 1,
		// "transmissoesRecebidas.eventoTransmissao":           1,
		// "transmissoesRecebidas.idPontoInteresse":            1,
		// "transmissoesRecebidas":                             1,
		// "lineString":                                        1,
		// "clienteId":                                         1,
		// "contador":                                          1,
		// "idRotaAberturaViagem":                              1,
		// "numeroLinhaArrastado":                              1,
		// "arrastoAutomatico":                                 1,
		// "dataFimAtraso":                                     1,
		// "kmPercurso":                                        1,
		// "codigoCobrador":                                    1,
		// "velocidadeMedia":                                   1,
		// "ipk":                                               1,
		// "tempoViagem":                                       1,
		// "diferencaPlanejado":                                1,
		// "passageiros":                                       1,
		// "descrIdRota":                                       1,
		// "excluido":                                          1,
		// "dataCriacao":                                       1,
		// "dataCriacaoRegistro":                               1,
		// "dtUltimaViagemAberta":                              1,
	}

	situacoes := [...]int{1, 2, 3, 4, 5, 7}
	// situacoes := [...]int{4}

	trajetos := make([]*bson.ObjectId, len(filtro.ListaTrajetos))
	for i, t := range filtro.ListaTrajetos {
		trajetos[i] = t.ID
	}
	dtInicio := filtro.GetDataInicio()
	dtFim := filtro.GetDataFim()

	query := bson.M{
		"excluido":                     false,
		"clienteId":                    filtro.IDCliente,
		"situacaoAtual":                bson.M{"$in": situacoes},
		"partida.trajetoExecutado._id": bson.M{"$in": trajetos},
		"executada.dataInicio":         bson.M{"$gte": dtInicio, "$lte": dtFim},
		// "partida":                      bson.M{"$exists": true},
		// "alocacao.idHorario":           bson.M{"$ne": ""},
		// "_id": bson.M{"$in": []bson.ObjectId{
		// 	bson.ObjectIdHex("5b6f38c3e4b0ad466e14ac3e"),
		// 	bson.ObjectIdHex("5b6f31bce4b0ad466e14ab7a")},
		// },
	}

	collection := session.DB(cfg.Config.MongoDB.Database).C("ViagemExecutada")
	var q *mgo.Query
	q = collection.Find(query)
	q.Select(projecao)
	// q.Limit(1)
	err = q.All(&listarViagens)

	logger.Debugf("viagemExecutada.size %d\n", len(listarViagens))
	/** /
	logger.Tracef("%#v\n", listarViagens)
	/**/

	if err != nil {
		logger.Errorf("Erro ao Listar Viagens no mongodb %s\n", err)
	}

	return listarViagens, err
}
