package webservice

import (
	"encoding/json"
	"net/http"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/intercept"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/service/viagemplanejamento"

	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

var logger logging.Logger

//InitConfig - é responsável por iniciar configuração da package
func InitConfig() {
	logger = logging.NewLogger("webservice", cfg.Config.Logging.Level)
}

var viagemplanejamentoService *viagemplanejamento.Service
var router *httprouter.Router

type myWeb struct {
}

func (c myWeb) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	res.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PUT, DELETE, OPTIONS, HEAD, PATCH")
	router.ServeHTTP(res, req)
	req.Body.Close()
}

//InitServer é responsável por inicializar o servidor http
func InitServer() {
	carragarDependencias()

	router = httprouter.New()

	mid := middleware.NewStack()

	mid.Use(intercept.ValidaToken)

	router.POST("/v1/viagemPlanejamento/filtrar", ConsultaViagemPlanejamento)
	router.POST("/api/v1/planejamentoviagem/dashboard", ConsultaViagemPlanejamento)

	logger.Infof("Servidor rodando na porta %v\n", cfg.Config.Server.Port)
	err := http.ListenAndServe(":"+cfg.Config.Server.Port, myWeb{})

	if err != nil {
		logger.Errorf("Erro ao subir o servidor na porta %v - %s\n", cfg.Config.Server.Port, err)
	}

}

func carragarDependencias() error {
	var err error
	con, err := database.GetSQLConnection()
	if err != nil {
		return err
	}
	planEscRep := repository.NewPlanejamentoEscalaRepository(con)

	mongoDB, err := database.GetMongoDB()
	if err != nil {
		return err
	}
	vigExecRep := repository.NewViagemExecutadaRepository(mongoDB)

	viagemplanejamentoService = viagemplanejamento.NewViagemPlanejamentoService(planEscRep, vigExecRep)
	return err
}

//ConsultaViagemPlanejamento - é responsável pela consulta de Viagens x Planejamento
func ConsultaViagemPlanejamento(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	var filter dto.FilterDTO
	err := json.NewDecoder(req.Body).Decode(&filter)
	if err != nil {
		logger.Errorf("Erro ao converter filtro %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Tracef("FILTRO: %#v\n", filter)

	consultaViagemPlanejamentoDTO, err := viagemplanejamentoService.Consultar(filter)

	if err != nil {
		logger.Errorf("Erro ConsultarViagemPlanejamento %+v - %s\n", filter, err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(consultaViagemPlanejamentoDTO)
}
