package viagemplanejamento

import (
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"gopkg.in/mgo.v2/bson"
)

//Cache -
var Cache = CacheViagemplanejamento{
	TrajetoLinha: make(map[*bson.ObjectId]dto.LinhaDTO),
}

//CacheViagemplanejamento -
type CacheViagemplanejamento struct {
	TrajetoLinha map[*bson.ObjectId]dto.LinhaDTO
}
