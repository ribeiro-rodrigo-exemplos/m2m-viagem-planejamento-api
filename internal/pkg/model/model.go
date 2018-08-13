package model

import (
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Modelo - MySQL Procedure

//ProcPlanejamentoEscala - Mapeamento de Stored Procedure
type ProcPlanejamentoEscala struct {
	IDPlanejamento int32     `json:"id_planejamento"`
	IDTrajeto      string    `json:"id_trajeto"`
	IDHorario      int32     `json:"id_horario"`
	IDTabela       int32     `json:"id_tabela"`
	NmTabela       string    `json:"nm_tabela"`
	IDEmpresaPlan  *int32    `json:"id_empresa_plan"`
	Partida        time.Time `json:"partida"`
	Chegada        time.Time `json:"chegada"`
	CodVeiculoPlan int32     `json:"cod_veiculo_plan"`
}

//Modelo - MySQL

//Cliente -
type Cliente struct {
	IDCliente int32  `json:""`
	Nome      string `json:""`
	Timezone  string `json:""`
	Location  *time.Location
}

// AtualizarLocation -
func (c *Cliente) AtualizarLocation() {
	loc, err := time.LoadLocation(c.Timezone)
	if err == nil {
		c.Location = loc
	}

}

//Modelo - MongoDB

//MensagemObs -
type MensagemObs struct {
	Mensagem        string    `bson:"mensagem"`
	DataAtualizacao time.Time `bson:"dataAtualizacao"`
	UsuarioCriacao  string    `bson:"usuarioCriacao"`
	Excluido        bool      `bson:"excluido"`
}

//GpsLinha -
type GpsLinha struct {
	Coordinates [][]float64 `bson:"coordinates"`
	TypeAtt     string      `bson:"type"`
}

//Transmissao -
type Transmissao struct {
	dataTransmissao time.Time
}

//TransmissoesRecebidas -
type TransmissoesRecebidas struct {
	IDPontoInteresse  string      `bson:"idPontoInteresse"`
	Transmissao       Transmissao `bson:"transmissao"`
	EventoTransmissao int32       `bson:"eventoTransmissao"`
	IDTrajeto         string      `bson:"idTrajeto"`
}

//Contador -
type Contador struct {
	Contador1 int32 `bson:"contador1"`
	Contador2 int32 `bson:"contador2"`
	Contador3 int32 `bson:"contador3"`
}

//Planejada -
type Planejada struct {
	DataFim    time.Time     `bson:"dataFim"`
	DataInicio time.Time     `bson:"dataInicio"`
	Veiculo    VeiculoViagem `bson:"veiculo"`
	TabelaID   string        `bson:"tabelaId"`
}

//VeiculoViagem -
type VeiculoViagem struct {
	ID        int32  `bson:"id"`
	Prefixo   string `bson:"prefixo"`
	IDEmpresa int32  `bson:"idEmpresa"`
}

//Executada -
type Executada struct {
	DataFim    *time.Time    `bson:"dataFim"`
	DataInicio time.Time     `bson:"dataInicio"`
	Veiculo    VeiculoViagem `bson:"veiculo"`
}

//Alocacao -
type Alocacao struct {
	IDPlanejamento string `bson:"idPlanejamento"`
	IDEscala       string `bson:"idEscala"`
	IDHorario      string `bson:"idHorario"`
}

//ViagemExecutada -
type ViagemExecutada struct {
	Alocacao              Alocacao                `bson:"alocacao"`
	ID                    bson.ObjectId           `bson:"_id"`
	ClienteID             int32                   `bson:"clienteId"`
	SituacaoAtual         int32                   `bson:"situacaoAtual"`
	Executada             Executada               `bson:"executada"`
	PorcentagemConclusao  string                  `bson:"porcentagemConclusao"`
	Planejada             Planejada               `bson:"planejada"`
	Contador              Contador                `bson:"contador"`
	TransmissoesRecebidas []TransmissoesRecebidas `bson:"transmissoesRecebidas"`
	LineString            GpsLinha                `bson:"lineString"`
	IDRotaAberturaViagem  string                  `bson:"idRotaAberturaViagem"`
	NumeroLinhaArrastado  string                  `bson:"numeroLinhaArrastado"`
	ArrastoAutomatico     bool                    `bson:"arrastoAutomatico"`
	TipoViagem            int32                   `bson:"tipoViagem"`
	DataFimAtraso         time.Time               `bson:"dataFimAtraso"`
	KmPercurso            float64                 `bson:"kmPercurso"`
	CodigoMotorista       string                  `bson:"codigoMotorista"`
	CodigoCobrador        string                  `bson:"codigoCobrador"`
	VelocidadeMedia       float64                 `bson:"velocidadeMedia"`
	Ipk                   float64                 `bson:"ipk"`
	TempoViagem           int64                   `bson:"tempoViagem"`
	DiferencaPlanejado    int32                   `bson:"diferencaPlanejado"`
	QntPassageiros        int32                   `bson:"qntPassageiros"`
	Passageiros           []time.Time             `bson:"passageiros"`
	DescrIDRota           string                  `bson:"descrIdRota"`
	Excluido              bool                    `bson:"excluido"`
	DataCriacao           time.Time               `bson:"dataCriacao"`
	DataCriacaoRegistro   time.Time               `bson:"dataCriacaoRegistro"`
	DtUltimaViagemAberta  time.Time               `bson:"dtUltimaViagemAberta"`
	MensagemObservacao    MensagemObs             `bson:"mensagemObs"`
	Partida               Partida                 `bson:"partida"`
}

//Partida -
type Partida struct {
	TrajetoExecutado TrajetoExecutado `bson:"trajetoExecutado"`
}

//TrajetoExecutado -
type TrajetoExecutado struct {
	IDObject bson.ObjectId `bson:"_id"`
}

//Modelo - API Planejamento

//Empresa -
type Empresa struct {
	EmpresaID int32 `json:"empresaId"`
	ID        int32 `json:"id"`
	IDEmpresa int32 `json:"idEmpresa"`
}

//GetIDEmpresa -
func (e Empresa) GetIDEmpresa() (idEmpresa int32) {
	if e.IDEmpresa > 0 {
		idEmpresa = e.IDEmpresa
	} else if e.ID > 0 {
		idEmpresa = e.ID
	} else {
		idEmpresa = e.EmpresaID
	}
	return
}

//Horarios -
type Horarios struct {
	Chegada                       string `json:"chegada"`
	ToleranciaAtraso              string `json:"toleranciaAtraso"`
	ToleranciaAtrasoChegada       int32  `json:"toleranciaAtrasoChegada"`
	ToleranciaAdiantamento        string `json:"toleranciaAdiantamento"`
	ToleranciaAdiantamentoChegada int32  `json:"toleranciaAdiantamentoChegada"`
	ToleranciaAtraso1             string `json:"toleranciaAtraso1"`
	ToleranciaAtrasoPartida       int32  `json:"toleranciaAtrasoPartida"`
	ToleranciaAdiantamento1       string `json:"toleranciaAdiantamento1"`
	ToleranciaAdiantamentoPartida int32  `json:"toleranciaAdiantamentoPartida"`
	Partida                       string `json:"partida"`
	IDHorario                     int32  `json:"idHorario"`
}

//GetChegada -
func (h Horarios) GetChegada() (chegada string) {
	if h.Chegada != "" {
		chegadaTrim := strings.Trim(h.Chegada, " ")
		if len(chegadaTrim) < 8 {
			chegada = chegadaTrim + ":00"
		}
	}
	return
}

//GetPartida -
func (h Horarios) GetPartida() (partida string) {
	if h.Partida != "" {
		partidaTrim := strings.Trim(h.Partida, " ")
		if len(partidaTrim) < 8 {
			partida = partidaTrim + ":00"
		}
	}
	return
}

//GetToleranciaAdiantamento -
func (h Horarios) GetToleranciaAdiantamento() (toleranciaAdiantamento int32) {
	if h.ToleranciaAdiantamentoPartida != 0 {
		toleranciaAdiantamento = h.ToleranciaAdiantamentoPartida
	}
	conv, _ := strconv.Atoi(h.ToleranciaAdiantamento)
	toleranciaAdiantamento = int32(conv)
	return
}

//GetToleranciaAtraso -
func (h Horarios) GetToleranciaAtraso() (toleranciaAtraso int32) {
	if h.ToleranciaAtrasoPartida != 0 {
		toleranciaAtraso = h.ToleranciaAtrasoPartida
	}
	conv, _ := strconv.Atoi(h.ToleranciaAtraso)
	toleranciaAtraso = int32(conv)
	return
}

//GetToleranciaAdiantamentoChegada -
func (h Horarios) GetToleranciaAdiantamentoChegada() (toleranciaAdiantamentoChegada int32) {
	if h.ToleranciaAdiantamentoChegada != 0 {
		toleranciaAdiantamentoChegada = h.ToleranciaAdiantamentoChegada
	}
	conv, _ := strconv.Atoi(h.ToleranciaAdiantamento1)
	toleranciaAdiantamentoChegada = int32(conv)
	return
}

//GetToleranciaAtrasoChegada -
func (h Horarios) GetToleranciaAtrasoChegada() (toleranciaAtrasoChegada int32) {
	if h.ToleranciaAtrasoChegada != 0 {
		toleranciaAtrasoChegada = h.ToleranciaAtrasoChegada
	}
	conv, _ := strconv.Atoi(h.ToleranciaAtraso1)
	toleranciaAtrasoChegada = int32(conv)
	return
}

//Tabela -
type Tabela struct {
	ID        string     `json:"_id"`
	TabelaID  string     `json:"tabelaId"`
	Horarios  []Horarios `json:"horario"`
	IDTabela  int32      `json:"idTabela"`
	CodTabela string     `json:"codTabela"`
	Empresa   Empresa    `json:"empresa"`
}

//GetTabelaID -
func (t Tabela) GetTabelaID() (tabelaID string) {
	if t.CodTabela != "" {
		tabelaID = t.CodTabela
	}
	tabelaID = t.TabelaID
	return
}

//GetID -
func (t Tabela) GetID() (id string) {
	if t.IDTabela != 0 {
		id = strconv.Itoa(int(t.IDTabela))
	}
	id = t.ID
	return
}

//PlanejamentoTrajeto -
type PlanejamentoTrajeto struct {
	TipoDia        string   `json:"tipoDia"`
	Tabelap        []Tabela `json:"tabela"`
	IDPlanejamento string   `json:"idPlanejamento"`
}

//PlanejamentoTrajetoEmbedded -
type PlanejamentoTrajetoEmbedded struct {
	Planejamentos []PlanejamentoTrajeto `json:"planejamentos"`
}

//PlanejamentoTrajetoAPIResponse -
type PlanejamentoTrajetoAPIResponse struct {
	Embedded PlanejamentoTrajetoEmbedded `json:"_embedded"`
}
