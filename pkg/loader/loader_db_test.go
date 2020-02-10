package loader

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var (
	Database *sql.DB
)

// Establish database connection
func setupDatabase() {
	const (
		host     = "localhost"
		port     = 5432
		user     = "root"
		password = "root"
		dbname   = "testdb"
	)

	connectionString := fmt.Sprintf("host=%s port=%d user=%s " +
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	Database, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	// Validate database connection
	err = Database.Ping()
	if err != nil {
		panic(err)
	}

	SeedDatabase()
}

// Close database connection
func teardownDatabase() {
	err := Database.Close()
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setupDatabase()
	m.Run()
	teardownDatabase()
}

func SeedDatabase() {
	seedData, err := ioutil.ReadFile("seedData.sql")
	if err != nil {
		panic(err)
	}

	_, err = Database.Exec(string(seedData))
}

func TestGetUserUuid(t *testing.T) {
	userUUID, domainName, err := GetUserUuid(Database, "28dc4965-8d0b-484d-bf8a-49986c53ef4e", "4d902414-82f5-427c-93df-bb3cb494756a")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	assert.Equal(t, userUUID, "3685ab1d-2e05-44fe-8c31-87805021e189")
	assert.Equal(t, domainName, "domain1")
}