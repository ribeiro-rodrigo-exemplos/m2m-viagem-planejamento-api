package webservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/cache"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/intercept"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/service/viagemplanejamento"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

var logger logging.Logger

//InitConfig - é responsável por iniciar configuração da package
func InitConfig() {
	logger = logging.NewLogger("webservice", cfg.Config.Logging.Level)
}

var viagemplanejamentoService chan *viagemplanejamento.Service
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

type request struct {
	filter     dto.FilterDTO
	res        *response
	chResponse chan *response
}

type response struct {
	dto *dto.ConsultaViagemPlanejamentoDTO
	err error
	req *request
}

//InitServer é responsável por inicializar o servidor http
func InitServer() {
	carragarDependencias()

	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	// defaultTransportPointer.MaxIdleConns = 100
	defaultTransportPointer.MaxIdleConnsPerHost = cfg.Config.HTTP.Transport.MaxIdleConnsPerHost
	logger.Infof("HTTP.Transport.MaxIdleConnsPerHost %v", cfg.Config.HTTP.Transport.MaxIdleConnsPerHost)

	logger.Infof("HTTP.Request.MaxConcurrent %v", cfg.Config.HTTP.Request.MaxConcurrent)

	router = httprouter.New()

	mid := middleware.NewStack()

	// mid.Use(intercept.ValidaToken)

	intercept.ConfigRateLimit()
	mid.Use(intercept.RateLimit)

	router.POST("/v1/viagemPlanejamento/filtrar", mid.Wrap(ConsultaViagemPlanejamento))
	router.POST("/api/v1/planejamentoviagem/dashboard", mid.Wrap(ConsultaViagemPlanejamento))
	router.PUT("/api/v1/planejamentoviagem/dashboard", mid.Wrap(ConsultaViagemPlanejamentoDashboard))

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

	logger.Infof("ViagemPlanejamento.MaxConcurrent %d ", cfg.Config.Service.ViagemPlanejamento.MaxConcurrent)

	viagemplanejamentoService = make(chan *viagemplanejamento.Service, cfg.Config.Service.ViagemPlanejamento.MaxConcurrent*2)
	for i := 0; i < cfg.Config.Service.ViagemPlanejamento.MaxConcurrent; i++ {
		viagemplanejamentoService <- viagemplanejamento.NewViagemPlanejamentoService(planEscRep, vigExecRep, cacheCliente)
	}
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

	vps := <-viagemplanejamentoService
	consultaViagemPlanejamentoDTO, err := vps.Consultar(filter)
	viagemplanejamentoService <- vps

	if err != nil {
		logger.Errorf("ConsultarViagemPlanejamento %s - %+v\n", err, filter)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode("Falha ao ConsultarViagemPlanejamento")
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

	listaTrajetos := make([]bson.ObjectId, len(filter.ListaTrajetos))
	for i := 0; i < len(filter.ListaTrajetos); i++ {
		t := filter.ListaTrajetos[i]
		listaTrajetos[i] = t.ID
	}

	filterAdaptado := dto.FilterDTO{
		ListaTrajetos: listaTrajetos,
		IDCliente:     filter.IDCliente,
		Ordenacao:     filter.Ordenacao,
		DataInicio:    filter.DataInicio + " " + strings.Replace(filter.HoraInicio, " ", "", -1),
		DataFim:       filter.DataFim + " " + strings.Replace(filter.HoraFim, " ", "", -1),
	}

	vps := <-viagemplanejamentoService
	consultaViagemPlanejamentoDTO, err := vps.Consultar(filterAdaptado)
	viagemplanejamentoService <- vps

	if err != nil {
		logger.Errorf("ConsultarViagemPlanejamento %s - %+v\n", err, filter)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode("Falha ao ConsultarViagemPlanejamento")
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(consultaViagemPlanejamentoDTO)
}
