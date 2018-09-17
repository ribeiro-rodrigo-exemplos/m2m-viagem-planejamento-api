package repository

import (
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//PontoInteresseRepository -
type PontoInteresseRepository struct {
	mongoDB *database.MongoDB
}

//NewPontoInteresseRepository -
func NewPontoInteresseRepository(mongoDB *database.MongoDB) *PontoInteresseRepository {
	p := new(PontoInteresseRepository)
	p.mongoDB = mongoDB
	return p
}

//ConsultarPorID -
func (p *PontoInteresseRepository) ConsultarPorID(id bson.ObjectId) (*model.PontoInteresse, error) {
	var pontoInteresse *model.PontoInteresse
	var err error

	session, err := p.mongoDB.GetSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	projecao := bson.M{
		"_id":  1,
		"nome": 1,
	}

	query := bson.M{
		"_id": id,
	}

	collection := session.DB(cfg.Config.MongoDB.Database).C("PontoInteresse")
	var q *mgo.Query
	q = collection.Find(query)
	q.Select(projecao)
	err = q.One(&pontoInteresse)

	return pontoInteresse, err
}

//CarregarMapaPontoInteresses -
func (p *PontoInteresseRepository) CarregarMapaPontoInteresses(ids []bson.ObjectId) (map[bson.ObjectId]*model.PontoInteresse, error) {
	pontoInteresse := make(map[bson.ObjectId]*model.PontoInteresse)
	var err error

	session, err := p.mongoDB.GetSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	projecao := bson.M{
		"_id":  1,
		"nome": 1,
	}

	query := bson.M{
		"_id": bson.M{"$in": ids},
	}

	collection := session.DB(cfg.Config.MongoDB.Database).C("PontoInteresse")
	var q *mgo.Query
	q = collection.Find(query)
	q.Select(projecao)
	iter := q.Iter()

	var pontoInteresseRetorno struct {
		ID   bson.ObjectId `bson:"_id"`
		Nome string        `bson:"nome"`
	}
	for iter.Next(&pontoInteresseRetorno) {
		pontoInteresse[pontoInteresseRetorno.ID] = &model.PontoInteresse{
			ID:   pontoInteresseRetorno.ID,
			Nome: pontoInteresseRetorno.Nome,
		}
	}
	if err = iter.Err(); err != nil {
		return nil, err
	}
	if err = iter.Close(); err != nil {
		return nil, err
	}

	return pontoInteresse, err
}

//ListarIdentificacaoPontosFinal -
func (p *PontoInteresseRepository) ListarIdentificacaoPontosFinal() ([]bson.ObjectId, error) {
	var ids []bson.ObjectId
	var err error

	session, err := p.mongoDB.GetSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	projecao := bson.M{
		"_id": 1,
		"trajetos.endPoint._id": 1,
	}

	query := bson.M{
		// "trajetos.endPoint._id": bson.M{"$exists": true},
	}

	collection := session.DB(cfg.Config.MongoDB.Database).C("Linha")
	var q *mgo.Query
	q = collection.Find(query)
	q.Select(projecao)
	iter := q.Iter()

	idsUnicos := make(map[string]interface{})
	var linha Linha

	for iter.Next(&linha) {
		for _, t := range linha.Trajetos {
			if !t.EndPoint.ID.Valid() {
				continue
			}
			id := t.EndPoint.ID.Hex()
			if _, k := idsUnicos[id]; !k {
				idsUnicos[id] = nil
				ids = append(ids, t.EndPoint.ID)
			}
		}
	}
	if err = iter.Err(); err != nil {
		return nil, err
	}
	if err = iter.Close(); err != nil {
		return nil, err
	}

	return ids, err
}

//Ponto -
type Ponto struct {
	ID bson.ObjectId `bson:"_id"`
}

//Trajeto -
type Trajeto struct {
	EndPoint Ponto `bson:"endPoint"`
}

//Linha -
type Linha struct {
	Trajetos []Trajeto `bson:"trajetos"`
}
