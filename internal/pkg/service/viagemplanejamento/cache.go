package viagemplanejamento

import (
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
)

//Cache -
var Cache = CacheViagemplanejamento{
	TrajetoLinha: make(map[string]dto.LinhaDTO),
}

//CacheViagemplanejamento -
type CacheViagemplanejamento struct {
	TrajetoLinha map[string]dto.LinhaDTO
}
