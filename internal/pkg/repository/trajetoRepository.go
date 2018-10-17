package repository

import (
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var projecao = bson.M{
	"_id":                   1,
	"descr":                 1,
	"numero":                1,
	"trajetos._id":          1,
	"trajetos.nome":         1,
	"trajetos.ativo":        1,
	"trajetos.endPoint._id": 1,
}

//TrajetoRepository -
type TrajetoRepository struct {
	mongoDB *database.MongoDB
}

//NewTrajetoRepository -
func NewTrajetoRepository(mongoDB *database.MongoDB) *TrajetoRepository {
	p := new(TrajetoRepository)
	p.mongoDB = mongoDB
	return p
}

//ConsultarPorID -
func (p *TrajetoRepository) ConsultarPorID(id bson.ObjectId) (model.Trajeto, error) {
	var trajeto model.Trajeto
	var err error

	session, err := p.mongoDB.GetSession()

	if err != nil {
		return trajeto, err
	}

	defer session.Close()

	query := bson.M{
		"trajetos._id": id,
	}

	var linha *model.Linha

	collection := session.DB(cfg.Config.MongoDB.Database).C("Linha")
	var q *mgo.Query
	q = collection.Find(query)
	q.Select(projecao)
	err = q.One(&linha)

	if err != nil {
		return trajeto, err
	}

	if linha != nil {
		for _, t := range linha.Trajetos {
			if t.ID == id {
				trajeto = t
				break
			}
		}
	}

	return trajeto, err
}

//CarregarMapaTrajetos -
func (p *TrajetoRepository) CarregarMapaTrajetos() (map[string]model.Trajeto, error) {
	mapaTrajetos := make(map[string]model.Trajeto)
	var err error

	session, err := p.mongoDB.GetSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	query := bson.M{}

	collection := session.DB(cfg.Config.MongoDB.Database).C("Linha")
	var q *mgo.Query
	q = collection.Find(query)
	q.Select(projecao)
	iter := q.Iter()

	var linha model.Linha

	for iter.Next(&linha) {
		for _, t := range linha.Trajetos {
			if t.Ativo {
				t.Linha = model.Linha{Numero: linha.Numero}
				mapaTrajetos[t.ID.Hex()] = t
			}
		}
	}
	if err = iter.Err(); err != nil {
		return nil, err
	}
	if err = iter.Close(); err != nil {
		return nil, err
	}
	return mapaTrajetos, err
}

//ListarTrajetos -
func (p *TrajetoRepository) ListarTrajetos() ([]model.Trajeto, error) {
	var listaTrajetos []model.Trajeto
	var err error

	session, err := p.mongoDB.GetSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	query := bson.M{}

	collection := session.DB(cfg.Config.MongoDB.Database).C("Linha")
	var q *mgo.Query
	q = collection.Find(query)
	q.Select(projecao)
	iter := q.Iter()

	var linha model.Linha

	for iter.Next(&linha) {
		nl := model.Linha{ID: bson.ObjectIdHex(linha.ID.Hex())}
		for _, t := range linha.Trajetos {
			t.Linha = nl
			listaTrajetos = append(listaTrajetos, t)
		}
	}
	if err = iter.Err(); err != nil {
		return nil, err
	}
	if err = iter.Close(); err != nil {
		return nil, err
	}

	return listaTrajetos, err
}
