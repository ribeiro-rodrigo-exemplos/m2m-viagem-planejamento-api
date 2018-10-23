package viagemplanejamento

import (
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
)

//CacheViagemplanejamento -
type CacheViagemplanejamento struct {
	TrajetoLinha map[string]dto.LinhaDTO
}
