package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pawlobanano/UGF3ZcWCIEdvZ29BcHBzIE5BU0E/config"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	err = testDB.Ping()
	if err != nil {
		log.Fatal("error pinging DB", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
