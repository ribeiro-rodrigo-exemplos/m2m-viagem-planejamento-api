package intercept

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//BodyLogger - logging Request.Body em nível Debug
func BodyLogger(next httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if loggerRequestBody.IsDebugEnabled() {
			buf, bodyErr := ioutil.ReadAll(req.Body)
			if bodyErr != nil {
				loggerRequestBody.Errorf("%s\n", bodyErr.Error())
				http.Error(res, "Falha ao recuperar Request.Body", http.StatusInternalServerError)
			} else {
				rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
				loggerRequestBody.Debugf("\n%v\n", rdr1)
				rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
				req.Body = rdr2
			}
		}
		next(res, req, params)
	}
}
