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
func (c *PlanejamentoEscalaRepository) ListarPlanejamentosEscala(filtro *dto.FilterDTO, cache map[int16]*model.Cliente) ([]*model.ProcPlanejamentoEscala, error) {
	planejamentosEscala := []*model.ProcPlanejamentoEscala{}
	var err error

	var (
		dataPlan       time.Time
		idPlanejamento int32
		idTrajeto      string
		idHorario      int32
		idTabela       int32
		nmTabela       string
		idEmpresaPlan  int32
		partida        types.RawTime
		chegada        types.RawTime
		codVeiculoPlan sql.NullInt64
	)

	var sql string
	// sql = "call sp_planejamento_vigente ('2018-07-24 00:00:00', '2018-07-24 23:59:59', '209', '\"O\",\"E\",\"3\",\"U\"', '\"555b6e830850536438063762\"') "
	sql = "call sp_planejamento_vigente (?, ?, ?, ?, ?) "
	logger.Tracef("{%s}\n", sql)

	dtInicio := filtro.GetDataInicioString()
	dtFim := filtro.GetDataFimString()
	idCliente := strconv.Itoa(filtro.IDCliente)
	var tiposDia = make([]string, len(filtro.TipoDia))
	for i, tipoDia := range filtro.TipoDia {
		tiposDia[i] = "\"" + tipoDia + "\""
	}
	var listaTrajetos = make([]string, len(filtro.ListaTrajetos))
	for i, trajeto := range filtro.ListaTrajetos {
		listaTrajetos[i] = "\"" + trajeto.Hex() + "\""
	}

	rows, err := c.connection.Query(sql, dtInicio, dtFim, idCliente, strings.Join(tiposDia, ","), strings.Join(listaTrajetos, ","))
	if err != nil {
		logger.Errorf("%s\n", err)
		return nil, err
	}
	defer rows.Close()
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
		timeChegada, _ := chegada.Time()
		if err != nil {
			return nil, fmt.Errorf("ListarPlanejamentosEscala - Recuperação Chegada: %s\n ", err)
		}
		var veiculo int32
		if codVeiculoPlan.Valid {
			veiculo = int32(codVeiculoPlan.Int64)
		}

		//TODO - Obter Location do cache do cliente
		loc, err := time.LoadLocation("America/Sao_Paulo")
		if err != nil {
			panic(err)
		}

		planejamentoEscala := &model.ProcPlanejamentoEscala{
			IDPlanejamento: idPlanejamento,
			IDTrajeto:      idTrajeto,
			IDHorario:      idHorario,
			IDTabela:       idTabela,
			NmTabela:       nmTabela,
			IDEmpresaPlan:  idEmpresaPlan,
			Partida:        util.Concatenar(dataPlan, timePartida, loc),
			Chegada:        util.Concatenar(dataPlan, timeChegada, loc),
			CodVeiculoPlan: veiculo,
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
