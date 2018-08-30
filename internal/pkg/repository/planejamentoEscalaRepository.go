package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/util"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/database/types"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
)

//PlanejamentoEscalaRepository -
type PlanejamentoEscalaRepository struct {
	connection *sql.DB
}

//NewPlanejamentoEscalaRepository -
func NewPlanejamentoEscalaRepository(connection *sql.DB) *PlanejamentoEscalaRepository {
	c := new(PlanejamentoEscalaRepository)
	c.connection = connection
	return c
}

//ListarPlanejamentosEscala -
func (c *PlanejamentoEscalaRepository) ListarPlanejamentosEscala(filtro *dto.FilterDTO) ([]*model.ProcPlanejamentoEscala, error) {
	planejamentosEscala := []*model.ProcPlanejamentoEscala{}
	var err error

	var (
		dataPlan                *time.Time
		idPlanejamento          *int32
		idTrajeto               *string
		idHorario               *int32
		idTabela                *int32
		nmTabela                *string
		idEmpresaPlan           *int32
		partida                 *types.RawTime
		chegada                 *types.RawTime
		codVeiculoPlan          *string
		toleranciaAtrasoPartida *int32
		mensagensObservacao     *string
	)

	var sql string
	// sql = "call sp_planejamento_vigente ('2018-07-24 00:00:00', '2018-07-24 23:59:59', '209', '\"O\",\"E\",\"3\",\"U\"', '\"555b6e830850536438063762\"') "
	sql = "call sp_planejamento_vigente (?, ?, ?, ?, ?) "

	dtInicio := filtro.GetDataInicioString()
	dtFim := filtro.GetDataFimString()
	idCliente := strconv.Itoa(int(filtro.IDCliente))
	var tiposDia = make([]string, len(filtro.TipoDia))
	for i, tipoDia := range filtro.TipoDia {
		tiposDia[i] = "\"" + tipoDia + "\""
	}
	var listaTrajetos = make([]string, len(filtro.ListaTrajetos))
	for i, trajeto := range filtro.ListaTrajetos {
		listaTrajetos[i] = "\"" + trajeto.ID.Hex() + "\""
	}

	if logger.IsDebugEnabled() {
		logger.Debugf("call sp_planejamento_vigente ('%s', '%s', '%s', '%s', '%s')\n", *dtInicio, *dtFim, idCliente, strings.Join(tiposDia, ","), strings.Join(listaTrajetos, ","))
	}

	rows, err := c.connection.Query(sql, dtInicio, dtFim, idCliente, strings.Join(tiposDia, ","), strings.Join(listaTrajetos, ","))
	if err != nil {
		logger.Errorf("%s\n", err)
		return nil, err
	}
	defer rows.Close()

	var loc *time.Location

	if filtro.Complemento.Cliente != nil {
		loc = filtro.Complemento.Cliente.Location
	}

	for rows.Next() {
		err := rows.Scan(
			&dataPlan,
			&idPlanejamento,
			&idTrajeto,
			&idHorario,
			&idTabela,
			&nmTabela,
			&idEmpresaPlan,
			&partida,
			&chegada,
			&codVeiculoPlan,
			&toleranciaAtrasoPartida,
			&mensagensObservacao,
		)

		if err != nil {
			logger.Errorf("%s\n", err)
			return nil, err
		}

		/** /
		logger.Debugf(
			"ProcPlanEsc: %v, %v, %v, %v, %v, %v, %v, %v, %v ",
			idPlanejamento, idTrajeto, idHorario, idTabela, nmTabela, idEmpresaPlan, partida, chegada, codVeiculoPlan,
		)
		/**/

		timePartida, err := partida.Time()
		if err != nil {
			return nil, fmt.Errorf("ListarPlanejamentosEscala - Recuperação Partida: %s\n ", err)
		}
		timeChegada, err := chegada.Time()
		if err != nil {
			return nil, fmt.Errorf("ListarPlanejamentosEscala - Recuperação Chegada: %s\n ", err)
		}

		partida := util.Concatenar(*dataPlan, timePartida, loc)
		chegada := util.Concatenar(*dataPlan, timeChegada, loc)

		observacoes, err := obterObservacoes(mensagensObservacao)
		if err != nil {
			return nil, fmt.Errorf("ListarPlanejamentosEscala - Recuperação Observações: %s\n ", err)
		}

		planejamentoEscala := &model.ProcPlanejamentoEscala{
			IDPlanejamento:          idPlanejamento,
			IDTrajeto:               idTrajeto,
			IDHorario:               idHorario,
			IDTabela:                idTabela,
			NmTabela:                nmTabela,
			IDEmpresaPlan:           idEmpresaPlan,
			Partida:                 &partida,
			Chegada:                 &chegada,
			CodVeiculoPlan:          codVeiculoPlan,
			ToleranciaAtrasoPartida: toleranciaAtrasoPartida,
			MensagensObservacao:     observacoes,
		}
		logger.Tracef("%#v\n", planejamentoEscala)
		planejamentosEscala = append(planejamentosEscala, planejamentoEscala)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	/**/
	logger.Debugf("planejamentosEscala.size %d\n", len(planejamentosEscala))
	/**/

	return planejamentosEscala, err
}

func obterObservacoes(mensagnesObservacao *string) ([]model.MensagemObservacaoProc, error) {
	mensagensObservacaoProc := []model.MensagemObservacaoProc{}
	var err error

	// mensagnesObservacaoAUX := "1|2018-08-28|3194235|Teste Vitor|1234|vitor.coelho"
	// mensagnesObservacao = &mensagnesObservacaoAUX

	if mensagnesObservacao == nil {
		return mensagensObservacaoProc, err
	}

	for _, mensagem := range strings.Split(*mensagnesObservacao, ";") {

		campos := strings.Split(mensagem, "|")
		if len(campos) == 6 {
			var id int
			var dataAtualizacao time.Time
			var idPlanejamento int
			var mensagem string
			var idUsuario int
			var nomeUsuario string
			// var excluido bool

			id, err = strconv.Atoi(campos[0])
			if err != nil {
				break
			}

			dataAtualizacao, err = util.ObterTimeDeAMD(campos[1])
			if err != nil {
				break
			}

			idPlanejamento, err = strconv.Atoi(campos[2])
			if err != nil {
				break
			}

			mensagem = campos[3]

			idUsuario, err = strconv.Atoi(campos[4])
			if err != nil {
				break
			}

			nomeUsuario = campos[5]

			msg := model.MensagemObservacaoProc{
				ID:              int32(id),
				IDPlanejamento:  int32(idPlanejamento),
				Mensagem:        mensagem,
				DataAtualizacao: dataAtualizacao,
				UsuarioCriacao: model.UsuarioProc{
					ID:   int32(idUsuario),
					Nome: nomeUsuario,
				},
			}
			mensagensObservacaoProc = append(mensagensObservacaoProc, msg)
		}
	}

	return mensagensObservacaoProc, err

}
