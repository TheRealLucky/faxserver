package main

import (
	"../../pkg/config"
	"../../pkg/loader"
	log "../../pkg/logger"
	"../../pkg/tiffer"
	"database/sql"
	"fmt"
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

	// testing loader package
	/*
	log.Info("loading default settings")
	default_settings, err := loader.Load_default_settings(Database)
	if err != nil {
		panic(err)
	}
	fmt.Println(default_settings)
	*/
	/*
	log.Info("loading account informations")
	account, err := loader.Get_account_informations(Database)
	if err != nil {
		panic(err)
	}
	fmt.Println(account)
	*/
	log.Info("getting user uuid")
	uuid, domain, err := loader.Get_user_uuid(Database, )

	tiffer.Create_folder("x","u")
	m := make(map[string]string)
	m["from"]="KundeA"
	m["to"]="KundeB"
	m["subject"]="xsx"
	file, err := tiffer.Create_pdf(m,"x","u")
	if err != nil {
		fmt.Println("failed")
	}
	fmt.Println(file)

}
