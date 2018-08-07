package intercept

import (
	"context"
	"net/http"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

//ValidaToken - é responsavel por validar o token e seu tempo de duração(expiração)
func ValidaToken(next httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

		token := req.Header.Get("Authorization")

		tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.Config.ChaveJWT), nil
		})

		claims, _ := tk.Claims.(jwt.MapClaims)
		idCliente := claims["clienteId"]

		var chave interface{}
		chave = "clienteID"

		ctx := context.WithValue(req.Context(), chave, idCliente)

		if err == nil && tk.Valid {
			next(res, req.WithContext(ctx), params)
			req.Body.Close()
		} else {
			logger.Warnf("Token inválido\n")
			res.WriteHeader(http.StatusForbidden)
			req.Body.Close()
		}

	}
}
