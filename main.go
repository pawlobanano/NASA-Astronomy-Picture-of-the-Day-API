package main

import (
	"log"

	"github.com/pawlobanano/UGF3ZcWCIEdvZ29BcHBzIE5BU0E/api"
	"github.com/pawlobanano/UGF3ZcWCIEdvZ29BcHBzIE5BU0E/config"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerPort)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
