package config

import (
	log "../logger"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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

// return a configuration struct, initialized with the values of the given config file.
func InitConfig(filename string) (*Configuration, error) {
	log.Info("initializing config")
	config := &Configuration{}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("failed to read yaml file, filename: %v", err)
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Error("failed to unmarshal yaml file: %v", err)
		return nil, err
	}

	if config.Interval < 120 {
		log.Fatal("interval has to be bigger than 120 seconds. ")
		return nil, errors.Errorf("interval has to be bigger than 120 seconds. interval is: %v", config.Interval)
	}

	log.Info("configuration initialized, configured values are:")
	log.Info("--- config.yaml ---")
	log.Info("interval: %v", config.Interval)
	log.Info("mailboxes: %v", config.Mailboxes)
	log.Info("debug: %v", config.Debug)
	log.Info("maxRoutines: %v", config.MaxRoutines)
	log.Info("database host: %v", config.Database.Host)
	log.Info("database port: %v", config.Database.Port)
	log.Info("database user: %v", config.Database.User)
	log.Info("database password: %v", config.Database.Password)
	log.Info("database name: %v", config.Database.DBName)
	log.Info("--- config.yaml ---")

	return config, nil
}

// returns a database, initialized with the values of the given configuration.
// furthermore opens a connection to this Database.
func InitDatabase(config *Configuration) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)
	db, err := connectDB(connString)

	if err != nil {
		log.Error("failed to connect to database: ", err)
	}
	return db, nil
}

func connectDB(connString string) (*sql.DB, error) {
	log.Info("connecting to database with connection string: %v", connString)

	db,err:= sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Errorf("failed to establish database connection: \n%v", err)
	} else {
		log.Info("database connected")
		return db, nil
	}

}
