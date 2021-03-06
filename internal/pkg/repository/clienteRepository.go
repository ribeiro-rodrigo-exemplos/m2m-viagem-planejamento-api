package repository

import (
	"database/sql"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/model"
)

//ClienteRepository -
type ClienteRepository struct {
	connection *sql.DB
}

//NewClienteRepository -
func NewClienteRepository(connection *sql.DB) *ClienteRepository {
	c := new(ClienteRepository)
	c.connection = connection
	return c
}

//CarregarMapaClientes -
func (c *ClienteRepository) CarregarMapaClientes() (map[int32]*model.Cliente, error) {
	mapaClientes := make(map[int32]*model.Cliente)
	var err error

	var (
		id       int32
		nome     string
		timezone string
	)
	rows, err := c.connection.Query("SELECT id_cliente, nome, ds_timezone from cliente WHERE 1=1 ")
	if err != nil {
		logger.Errorf("%s\n", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &nome, &timezone)
		if err != nil {
			logger.Errorf("%s\n", err)
			return nil, err
		}

		/** /
		logger.Debugf("%v, %v, %v", id, nome, timezone)
		/**/

		cliente := &model.Cliente{
			IDCliente: id,
			Nome:      nome,
			Timezone:  timezone,
		}
		cliente.AtualizarLocation()
		mapaClientes[id] = cliente

		logger.Tracef("%#v\n", cliente)
	}
	err = rows.Err()
	if err != nil {
		logger.Errorf("%s\n", err)
	}

	return mapaClientes, err
}
