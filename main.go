package main

import (
	"log"

	"github.com/EUDAT-GEF/Bridgit/api"
	"github.com/EUDAT-GEF/Bridgit/utils"
)

func main() {
	config, err := utils.ReadConfigFile("./def/config.json")
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	app := api.NewApp(config)
	err = app.Start()
	if err != nil {
		log.Fatal("Failed to start the service")
	}
}
