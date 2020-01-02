package main

import (
	"../../pkg/config"
	log "../../pkg/logger"
)

func main() {
	log.Info("starting application")
	config.InitConfig("conf.yaml")
}
