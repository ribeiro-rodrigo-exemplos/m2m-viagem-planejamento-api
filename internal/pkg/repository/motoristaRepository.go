package repository

import (
	"database/sql"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
)

const sqlCarregarMapaMotoristas = "" +
	"SELECT " +
	"	id_funcionario, nm_nomeFuncionario, nm_matricula " +
	"FROM " +
	"	funcionario " +
	"WHERE 1=1 " +
	"	AND id_tipoFuncionario = 1 AND fl_situacao = true " +
	"GROUP BY " +
	"	nm_matricula"

//MotoristaRepository -
type MotoristaRepository struct {
	connection *sql.DB
}

//NewMotoristaRepository -
func NewMotoristaRepository(connection *sql.DB) *MotoristaRepository {
	c := new(MotoristaRepository)
	c.connection = connection
	return c
}

//CarregarMapaMotoristas -
func (c *MotoristaRepository) CarregarMapaMotoristas() (map[string]*model.Motorista, error) {
	mapaMotoristas := make(map[string]*model.Motorista)
	var err error
	var (
		id        int32
		nome      string
		matricula string
	)

	rows, err := c.connection.Query(sqlCarregarMapaMotoristas)
	if err != nil {
		logger.Errorf("%s\n", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &nome, &matricula)
		if err != nil {
			logger.Errorf("%s\n", err)
			return nil, err
		}

		/** /
		logger.Debugf("%v, %v, %v", id, nome, matricula)
		/**/

		motorista := model.NewMotorista(id, nome, matricula)
		mapaMotoristas[matricula] = motorista

		logger.Tracef("%#v\n", motorista)
	}
	err = rows.Err()
	if err != nil {
		logger.Errorf("%s\n", err)
	}

	return mapaMotoristas, err
}
