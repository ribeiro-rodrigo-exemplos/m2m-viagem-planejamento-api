package repository

import (
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var projecaoLinha = bson.M{
	"_id":                   1,
	"descr":                 1,
	"numero":                1,
	"trajetos._id":          1,
	"trajetos.nome":         1,
	"trajetos.ativo":        1,
	"trajetos.endPoint._id": 1,
	"consorcio.consorcioId": 1,
}

//LinhaRepository -
type LinhaRepository struct {
	mongoDB *database.MongoDB
}

//NewLinhaRepository -
func NewLinhaRepository(mongoDB *database.MongoDB) *LinhaRepository {
	p := new(LinhaRepository)
	p.mongoDB = mongoDB
	return p
}

//Listar -
func (p *LinhaRepository) Listar() (linhas []model.Linha, err error) {

	session, err := p.mongoDB.GetSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	query := bson.M{
		//"excluido" : false,
	}

	collection := session.DB(cfg.Config.MongoDB.Database).C("Linha")
	var q *mgo.Query
	q = collection.Find(query)
	q.Select(projecaoLinha)

	iter := q.Iter()
	defer iter.Close()

	var linha model.Linha

	for iter.Next(&linha) {
		linhas = append(linhas, linha)
	}
	if err = iter.Err(); err != nil {
		return nil, err
	}

	return
}
