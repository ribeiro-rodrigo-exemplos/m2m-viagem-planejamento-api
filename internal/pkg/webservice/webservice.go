package webservice

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/cache"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/service/viagemplanejamento"

	"github.com/julienschmidt/httprouter"
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
	// req.Body.Close()
}

var requestQueue = make(chan request, 2)
var responseQueue = make(chan chan response, 2)

type response struct {
	dto *dto.ConsultaViagemPlanejamentoDTO
	err error
}

type request struct {
	filter   dto.FilterDTO
	response chan response
}

func consultarViagemPlanejamento(filtro dto.FilterDTO) (*dto.ConsultaViagemPlanejamentoDTO, error) {

	chRes := <-responseQueue

	consultaViagemPlanejamentoDTO, err := viagemplanejamentoService.Consultar(filtro)

	res := <-chRes

	return res.dto, res.err
	// return consultaViagemPlanejamentoDTO, err
}

//InitServer é responsável por inicializar o servidor http
func InitServer() {
	carragarDependencias()

	responseQueue <- make(chan response, 2)

	var process = func() {
		for r := range requestQueue {
			// router.ServeHTTP(*r.res, r.req)
			// f := r.next
			// f(*r.res, r.req, *r.params)
			// r.next(r.res, r.req, r.params)
		}
	}
	go process()

	var delegate = func() {
		for r := range requestQueue {
			consultaViagemPlanejamentoDTO, err := viagemplanejamentoService.Consultar(filtro)
			// router.ServeHTTP(*r.res, r.req)
			// f := r.next
			// f(*r.res, r.req, *r.params)
			// r.next(r.res, r.req, r.params)
		}
	}
	go delegate()

	router = httprouter.New()

	// mid := middleware.NewStack()

	// mid.Use(intercept.ValidaToken)

	// intercept.ConfigRateLimit()
	// mid.Use(intercept.RateLimit)

	// router.POST("/v1/viagemPlanejamento/filtrar", mid.Wrap(ConsultaViagemPlanejamento))
	router.POST("/v1/viagemPlanejamento/filtrar", ConsultaViagemPlanejamento)
	// router.POST("/api/v1/planejamentoviagem/dashboard", ConsultaViagemPlanejamento)
	// router.PUT("/api/v1/planejamentoviagem/dashboard", ConsultaViagemPlanejamentoDashboard)

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

	cacheCliente, err := cache.GetCliente(nil)
	if err != nil {
		return err
	}

	viagemplanejamentoService = viagemplanejamento.NewViagemPlanejamentoService(planEscRep, vigExecRep, cacheCliente)
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

	consultaViagemPlanejamentoDTO, err := consultarViagemPlanejamento(filter)

	if err != nil {
		logger.Errorf("Erro ConsultarViagemPlanejamento %+v - %s\n", filter, err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(consultaViagemPlanejamentoDTO)
}

//ConsultaViagemPlanejamentoDashboard - é responsável pela consulta de Viagens x Planejamento
func ConsultaViagemPlanejamentoDashboard(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	var filter dto.FilterDashboardDTO
	err := json.NewDecoder(req.Body).Decode(&filter)
	if err != nil {
		logger.Errorf("Erro ao converter filtro %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Tracef("FILTRO: %#v\n", filter)

	filterAdaptado := dto.FilterDTO{
		ListaTrajetos: filter.ListaTrajetos,
		IDCliente:     filter.IDCliente,
		Ordenacao:     filter.Ordenacao,
		DataInicio:    filter.DataInicio + " " + strings.Replace(filter.HoraInicio, " ", "", -1),
		DataFim:       filter.DataFim + " " + strings.Replace(filter.HoraFim, " ", "", -1),
	}

	consultaViagemPlanejamentoDTO, err := consultarViagemPlanejamento(filterAdaptado)

	if err != nil {
		logger.Errorf("Erro ConsultarViagemPlanejamento %+v - %s\n", filter, err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(consultaViagemPlanejamentoDTO)
}
