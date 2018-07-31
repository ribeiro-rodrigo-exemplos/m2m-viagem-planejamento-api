package database

import (
	"database/sql"
	"fmt"
	"time"

	cfg "git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/config"

	_ "github.com/go-sql-driver/mysql" //Carrega driver MYSQL
)

var connection *sql.DB

//GetSQLConnection -
func GetSQLConnection() (*sql.DB, error) {
	var err error
	if connection == nil {
		for index := 0; index < cfg.Config.MySQL.Reconnect; index++ {
			connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Config.MySQL.User, cfg.Config.MySQL.Password, cfg.Config.MySQL.Host, cfg.Config.MySQL.Port, cfg.Config.MySQL.Database)
			logger.Debugf("%s", connectionString)
			db, err := sql.Open("mysql",
				connectionString)

			if err != nil {
				logger.Errorf("Reconectar banco de dados devido a falha - %s\n", err)
				time.Sleep(cfg.Config.MySQL.ReconnectSleep * time.Second)
				continue
			}

			db.SetMaxIdleConns(cfg.Config.MySQL.MaxIdleConns)
			db.SetMaxOpenConns(cfg.Config.MySQL.MaxOpenConns)

			err = db.Ping()
			if err != nil {
				logger.Errorf("Reconectar banco de dados devido a falha - %s\n", err)
				time.Sleep(cfg.Config.MySQL.ReconnectSleep * time.Second)
				continue
			}
			connection = db
			break
		}
	}
	return connection, err
}
