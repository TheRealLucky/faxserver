package main

import (
	"../../pkg/config"
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
	fmt.Println(Database)
	// testing purposes
	//default_settings, err := loader.Load_default_settings(Database)
	if err != nil {
		panic(err)
	}
	//fmt.Println(default_settings)

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


	//testcode on linux
	/*
	ka := make(map[string][]string)
	t := make([]string,1)
	l := make([]string,1)
	t = append(t,"d8a9bea9-640d-423e-5de0-f6f517d1d529.pdf")
	t = append(t, "f3d72fa4-9584-486c-4877-133980b81860.pdf")
	l = append(l, "dummy.pdf","sample.pdf")
	ka["./debug/mailboxes/x/u/2020-01-02"]= t
	//ka["./testlogos/"] = l
	_, err = tiffer.Merge_pdf(ka)
	if err != nil {   fmt.Println(err)}*/

}