package config

import (
	log "../logger"
	"database/sql"
)

type Config struct {

}

var Database *sql.DB

func InitConfig() {
	log.Info("initializing config")
}
