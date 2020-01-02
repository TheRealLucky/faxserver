package main

import (
	"../../pkg/config"
	log "../../pkg/logger"
	"database/sql"
)

var Config *config.Configuration
var Database *sql.DB
import (
	"fmt"
	tiffer "../../pkg/tiffer"
)

func main() {
	fmt.Println("Hello World!")
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
