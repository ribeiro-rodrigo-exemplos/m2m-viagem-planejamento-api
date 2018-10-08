package webservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/cache"
	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/intercept"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/repository"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/service/viagemplanejamento"

	"github.com/NYTimes/gziphandler"
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

func serveHTTP(res http.ResponseWriter, req *http.Request) {
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

	logger.Infof("HTTP.Response.Gzip.Enable %v", cfg.Config.HTTP.Response.Gzip.Enable)

	router = httprouter.New()

	mid := middleware.NewStack()

	// mid.Use(intercept.ValidaToken)

	intercept.ConfigRateLimit()
	mid.Use(intercept.RateLimit)
	mid.Use(intercept.BodyLogger)

	router.POST("/v1/viagemPlanejamento/filtrar", mid.Wrap(ConsultaViagemPlanejamento))
	router.POST("/api/v1/planejamentoviagem/dashboard", mid.Wrap(ConsultaViagemPlanejamento))
	router.PUT("/api/v1/planejamentoviagem/dashboard", mid.Wrap(ConsultaViagemPlanejamentoDashboard))

	var defaultHandler http.Handler

	defaultHandler = http.HandlerFunc(serveHTTP)
	if cfg.Config.HTTP.Response.Gzip.Enable {
		defaultHandler = gziphandler.GzipHandler(defaultHandler)
	}

	http.Handle("/", defaultHandler)

	logger.Infof("Servidor rodando na porta %v\n", cfg.Config.Server.Port)
	err := http.ListenAndServe(":"+cfg.Config.Server.Port, nil)

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
	cacheMotorista, err := cache.GetMotorista(nil)
	if err != nil {
		return err
	}
	cacheTrajeto, err := cache.GetTrajeto(nil)
	if err != nil {
		return err
	}
	cachePontoInteresse, err := cache.GetPontoInteresse(nil)
	if err != nil {
		return err
	}
	cacheAgrupamento, err := cache.GetAgrupamento(nil)
	if err != nil {
		return err
	}

	logger.Infof("ViagemPlanejamento.MaxConcurrent %d ", cfg.Config.Service.ViagemPlanejamento.MaxConcurrent)

	viagemplanejamentoService = make(chan *viagemplanejamento.Service, cfg.Config.Service.ViagemPlanejamento.MaxConcurrent*2)
	for i := 0; i < cfg.Config.Service.ViagemPlanejamento.MaxConcurrent; i++ {
		viagemplanejamentoService <- viagemplanejamento.NewViagemPlanejamentoService(planEscRep, vigExecRep, cacheCliente, cacheMotorista, cacheTrajeto, cachePontoInteresse, cacheAgrupamento)
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

	if err != nil {
		logger.Errorf("ConsultarViagemPlanejamento %s - %+v\n", err, filter)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode("Falha ao ConsultarViagemPlanejamento")
		viagemplanejamentoService <- vps
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(consultaViagemPlanejamentoDTO)
	viagemplanejamentoService <- vps
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

	listaAgrupamentos := make([]dto.AgrupamentoDTO, len(filter.ListaAgrupamentos))
	for i := 0; i < len(filter.ListaAgrupamentos); i++ {
		agrupamentoID := filter.ListaAgrupamentos[i]
		grupoID, err := strconv.Atoi(agrupamentoID)
		if err != nil {
			logger.Errorf("Erro ao converter filtro %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		listaAgrupamentos[i] = dto.AgrupamentoDTO{ID: int32(grupoID)}
	}

	listaTrajetos := make([]dto.TrajetoDTO, len(filter.ListaTrajetos))
	for i := 0; i < len(filter.ListaTrajetos); i++ {
		t := filter.ListaTrajetos[i]
		listaTrajetos[i] = dto.TrajetoDTO{ID: t.ID, Descricao: t.Descricao, Sentido: t.Sentido, Linha: dto.LinhaDTO{Numero: t.NumeroLinha}}
	}

	listaEmpresas := make([]dto.EmpresaDTO, len(filter.ListaEmpresas))
	for i := 0; i < len(filter.ListaEmpresas); i++ {
		empresaID := filter.ListaEmpresas[i]
		listaEmpresas[i] = dto.EmpresaDTO{ID: empresaID}
	}

	dataInicio := filter.DataInicio + " " + strings.Replace(filter.HoraInicio, " ", "", -1)
	dataFim := filter.DataFim + " " + strings.Replace(filter.HoraFim, " ", "", -1)

	filterAdaptado := dto.FilterDTO{
		ListaAgrupamentos: listaAgrupamentos,
		ListaTrajetos:     listaTrajetos,
		ListaEmpresas:     listaEmpresas,
		IDCliente:         filter.IDCliente,
		Ordenacao:         filter.Ordenacao,
		DataInicio:        &dataInicio,
		DataFim:           &dataFim,
	}

	vps := <-viagemplanejamentoService
	consultaViagemPlanejamentoDTO, err := vps.ConsultarDashboard(filterAdaptado)

	if err != nil {
		logger.Errorf("ConsultarViagemPlanejamento %s - %+v\n", err, filter)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode("Falha ao ConsultarViagemPlanejamento")
		viagemplanejamentoService <- vps
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)

	json.NewEncoder(res).Encode(consultaViagemPlanejamentoDTO)
	viagemplanejamentoService <- vps
}
