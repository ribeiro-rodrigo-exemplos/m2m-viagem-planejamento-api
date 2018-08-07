package intercept

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var requestQueue = make(chan request)

//ConfigRateLimit -
func ConfigRateLimit() {

	var delegate = func() {
		for r := range requestQueue {
			// f := r.next
			// f(*r.res, r.req, *r.params)
			r.next(r.res, r.req, r.params)
		}
	}
	go delegate()

}

/** /
//RateLimit - é responsavel por limitar quantidade de requests tratadas simultaneamente
func RateLimit(next httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		requestQueue <- request{&res, req, &params, next}
		// next(res, req, params)
	}
}

type request struct {
	res    *http.ResponseWriter
	req    *http.Request
	params *httprouter.Params
	next   httprouter.Handle
}
/**/

//RateLimit - é responsavel por limitar quantidade de requests tratadas simultaneamente
func RateLimit(next httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		requestQueue <- request{res, req, params, next}
		// next(res, req, params)
	}
}

type request struct {
	res    http.ResponseWriter
	req    *http.Request
	params httprouter.Params
	next   httprouter.Handle
}
