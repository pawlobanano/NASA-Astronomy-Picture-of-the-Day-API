package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pawlobanano/UGF3ZcWCIEdvZ29BcHBzIE5BU0E/config"
	db "github.com/pawlobanano/UGF3ZcWCIEdvZ29BcHBzIE5BU0E/db/sqlc"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	runDbMigration(config.MigrationURL, config.DBSource)

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("error pinging DB", err)
	}

	_ = db.NewRepository(dbConn)

	// TODO Implement API + Server
}

func runDbMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up", err)
	}

	log.Println("db migrated successfully")
}
