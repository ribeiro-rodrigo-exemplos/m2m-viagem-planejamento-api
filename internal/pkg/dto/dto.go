package dto

import (
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/util"
	"gopkg.in/mgo.v2/bson"
)

//FilterDTO filtro para consultas
type FilterDTO struct {
	ListaTrajetos []bson.ObjectId `json:"lista_trajetos"`
	IDCliente     int32           `json:"id_cliente"`
	IDVeiculo     int             `json:"id_veiculo"`
	Ordenacao     []string        `json:"ordenacao"`
	DataInicio    string          `json:"data_inicio"`
	DataFim       string          `json:"data_fim"`
	TipoDia       []string
	Complemento   DadosComplementares
}

//FilterDashboardDTO filtro para consultas dashboard
type FilterDashboardDTO struct {
	ListaLinhas   []bson.ObjectId `json:"linhas"`
	ListaTrajetos []bson.ObjectId `json:"trajetos"`
	IDCliente     int32           `json:"idCliente"`
	Status        []string        `json:"status"`
	Ordenacao     []string        `json:"ordenacao"`
	DataInicio    string          `json:"dataInicio"`
	HoraInicio    string          `json:"horaInicio"`
	DataFim       string          `json:"dataFim"`
	HoraFim       string          `json:"horaFim"`
	Timezone      string          `json:"timezone"`
}

//GetDataInicio -
func (f *FilterDTO) GetDataInicio() time.Time {
	var dt time.Time
	dt, _ = util.ObterUTCTime(f.DataInicio)
	return dt
}

//GetDataInicioString -
func (f *FilterDTO) GetDataInicioString() string {
	var dt time.Time
	dt, err := util.ObterUTCTime(f.DataInicio)
	if err != nil {
		return ""
	}
	str := util.FormatarAMDHMS(dt)
	return str
}

//DadosComplementares -
type DadosComplementares struct {
	Cliente *model.Cliente
}

//GetDataFim -
func (f *FilterDTO) GetDataFim() time.Time {
	var dt time.Time
	dt, _ = util.ObterUTCTime(f.DataFim)
	return dt
}

//GetDataFimString -
func (f *FilterDTO) GetDataFimString() string {
	var dt time.Time
	dt, err := util.ObterUTCTime(f.DataFim)
	if err != nil {
		return ""
	}
	str := util.FormatarAMDHMS(dt)
	return str
}

//ViagemDTO estrutura usada para mapear dados enviados para a tela
type ViagemDTO struct {
	ID                  bson.ObjectId `json:"id"`
	IDViagemExecutada   bson.ObjectId `json:"idViagemExecutada"`
	IDTabela            int32         `json:"idTabela"`
	IDHorario           int32         `json:"idHorario"`
	IDEmpresa           int32         `json:"idEmpresa"`
	IDEmpresaPlanejada  int32         `json:"idEmpresaPlanejada"`
	Status              int           `json:"status"`
	EmExecucao          bool          `json:"emExecucao"`
	VeiculoPlan         string        `json:"veiculoPlan"`
	VeiculoReal         string        `json:"veiculoReal"`
	NmTabela            string        `json:"nmTabela"`
	PartidaOrdenacao    time.Time     `json:"-"`
	PartidaPlanTime     time.Time     `json:"-"`
	PartidaPlan         string        `json:"partidaPlan"`
	ChegadaPlanTime     time.Time     `json:"-"`
	ChegadaPlan         string        `json:"chegadaPlan"`
	DiffPartida         int64         `json:"diffPartida"`
	DiffPartidaStr      string        `json:"diffPartidaStr"`
	EntrouEmPlaca       string        `json:"entrouEmPlaca"`
	PartidaReal         string        `json:"partidaReal"`
	ChegadaReal         string        `json:"chegadaReal"`
	DiffChegada         int64         `json:"diffChegada"`
	DiffChegadaStr      string        `json:"diffChegadaStr"`
	QtdePassageiros     int32         `json:"qtdePassageiros"`
	Proxima             string        `json:"proxima"`
	PercentualConclusao string        `json:"percentualConclusao"`
	Editada             bool          `json:"editada"`
	Headway             int64         `json:"headway"`
	HeadwayStr          string        `json:"headwayStr"`
	Data                string        `json:"data"`
	DataAbertura        string        `json:"dataAbertura"`
	DataFechamento      string        `json:"dataFechamento"`
	Ipk                 float64       `json:"ipk"`
	CdMotorista         string        `json:"cdMotorista"`
	NmMotorista         string        `json:"nmMotorista"`
	Duracao             string        `json:"duracao"`
	DuracaoSeg          int64         `json:"duracaoSeg"`
	IDVeiculo           string        `json:"idVeiculo"`
	Motivo              string        `json:"motivo"`
	VelocidadeMedia     float64       `json:"velocidadeMedia"`
	CodTransportadora   int           `json:"codTransportadora"`
	TipoFrota           int           `json:"tipoFrota"`
	Cobrador            int           `json:"cobrador"`
	DistanciaPercorrida float32       `json:"distanciaPercorrida"`
	Placa               string        `json:"placa"`
}

//TotalizadoresDTO -
type TotalizadoresDTO struct {
	Planejadas           int32   `json:"planejadas"`
	PlanejadasAteMomento int32   `json:"planejadasAteMomento"`
	Realizadas           int32   `json:"realizadas"`
	RealizadasPlanejadas int32   `json:"realizadasPlanejadas"`
	EmAndamento          int32   `json:"emAndamento"`
	Canceladas           int32   `json:"canceladas"`
	Passageiros          int32   `json:"passageiros"`
	NaoIniciadas         int32   `json:"naoIniciadas"`
	NaoRealizadas        int32   `json:"naoRealizadas"`
	Reforco              int32   `json:"reforco"`
	Atrasada             int32   `json:"atrasada"`
	IndiceExecucao       []int32 `json:"indiceExecucao"`
	IndicePartida        []int32 `json:"indicePartida"`
}

//ConsultaViagemPlanejamentoDTO Mapeia resultado da consulta enviados para tela
type ConsultaViagemPlanejamentoDTO struct {
	Informacoes      map[string]interface{}   `json:"informacoes"`
	ViagensExecutada []*model.ViagemExecutada `json:"viagensExecutada"`
	Totalizadores    *TotalizadoresDTO        `json:"totalizadores"`
	Viagens          []*ViagemDTO             `json:"viagens"`
}
