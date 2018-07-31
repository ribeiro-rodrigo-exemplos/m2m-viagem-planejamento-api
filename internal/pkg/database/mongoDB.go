package database

import (
	"fmt"
	"time"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"

	mgo "gopkg.in/mgo.v2"
)

var mongoDB *MongoDB

//MongoDB -
type MongoDB struct {
	mongoSession *mgo.Session
	inicializado bool
}

//GetMongoDB -
func GetMongoDB() (*MongoDB, error) {
	var err error
	if mongoDB == nil {
		mongoDB = new(MongoDB)
		mongoDB.inicializado = true
	}
	return mongoDB, err
}

//GetSession - realiza a conexão com o MongoDB
func (m *MongoDB) GetSession() (*mgo.Session, error) {
	var err error
	var sessionCopy *mgo.Session

	if !m.inicializado {
		return nil, fmt.Errorf("Instância de MongoDB não iniciada corretamente")
	}

	if m.mongoSession == nil {
		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:    []string{cfg.Config.MongoDB.Host},
			Timeout:  cfg.Config.MongoDB.Timeout * time.Second,
			Database: cfg.Config.MongoDB.Database,
		}

		session, err := mgo.DialWithInfo(mongoDBDialInfo)

		if err != nil {
			logger.Errorf("Erro ao conectar com o mongo %s\n", err)
			return sessionCopy, err
		}
		m.mongoSession = session

		logger.Infof("Conectado ao mongo\n")
	}
	sessionCopy = m.mongoSession.Copy()
	return sessionCopy, err
}
