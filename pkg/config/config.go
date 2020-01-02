package config

import (
	log "../logger"
	"database/sql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	Interval 	int		`yaml:"interval"`
	Mailboxes 	int		`yaml:"mailboxes"`
	Debug		bool	`yaml:"debug"`
	MaxRoutines int		`yaml:"maxRoutines"`
	Database struct {
		Host 		string `yaml:"host"`
		Port 		string `yaml:"port"`
		User 		string `yaml:"user"`
		Password 	string `yaml:"password"`
		DBName 		string `yaml:"dbname"`
	}
}

var Database *sql.DB
var Config *Configuration

func InitConfig(filename string) {
	log.Info("initializing config")
	Config = &Configuration{}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("failed to read yaml file, filename: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, Config)
	if err != nil {
		log.Error("failed to unmarshal yaml file: %v", err)
	}

	log.Info("configuration initialized, configured values are:")
	log.Info("--- config.yaml ---")
	log.Info("interval: %v", Config.Interval)
	log.Info("mailboxes: %v", Config.Mailboxes)
	log.Info("debug: %v", Config.Debug)
	log.Info("maxRoutines: %v", Config.MaxRoutines)
	log.Info("database host: %v", Config.Database.Host)
	log.Info("database port: %v", Config.Database.Port)
	log.Info("database user: %v", Config.Database.User)
	log.Info("database password: %v", Config.Database.Password)
	log.Info("database name: %v", Config.Database.DBName)
	log.Info("--- config.yaml ---")
}
