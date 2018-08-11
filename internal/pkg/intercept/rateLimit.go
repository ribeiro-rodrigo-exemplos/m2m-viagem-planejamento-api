package intercept

import (
	"net/http"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"github.com/julienschmidt/httprouter"
)

var requestQueue chan struct{}

//ConfigRateLimit -
func ConfigRateLimit() {
	requestQueue = make(chan struct{}, cfg.Config.HTTP.Request.MaxConcurrent*2)

	for i := 0; i < cfg.Config.HTTP.Request.MaxConcurrent; i++ {
		requestQueue <- struct{}{}

	}
}

//RateLimit - é responsavel por limitar quantidade de requests tratadas simultaneamente
func RateLimit(next httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		r := <-requestQueue
		next(res, req, params)
		requestQueue <- r
	}
}
