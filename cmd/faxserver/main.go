package main

import (
	"../../pkg/config"
	log "../../pkg/logger"
	"database/sql"
)

var Config *config.Configuration
var Database *sql.DB

func main() {
	log.Info("starting application")
	Config, err := config.InitConfig("conf.yaml")
	if err != nil {
		panic(err)
	}

	Database, err := config.InitDatabase(Config)
	if err != nil {
		panic(err)
	}

	// testing purposes
	err = Database.Ping()
	if err != nil {
		panic(err)
	}
}
