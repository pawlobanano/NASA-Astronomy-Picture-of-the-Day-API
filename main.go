package main

import (
	"embed"
	"log"

	"github.com/pawlobanano/NASA-Astronomy-Picture-of-the-Day-API/api"
	"github.com/pawlobanano/NASA-Astronomy-Picture-of-the-Day-API/config"
)

//go:embed test/json
var test embed.FS

func main() {
	api.TestDir = test

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
