package main

import (
	"log"


	"github.com/EUDAT-GEF/Bridgit/utils"
	"github.com/EUDAT-GEF/Bridgit/api"
)

const configFilePath = "config.json"


func main() {
	config, err := utils.ReadConfigFile("./config/config.json")
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	app := api.NewApp(config)
	app.Start()


	//startHttpServer(Config.PortNumber, NewRouter())
	//
	//log.Fatal(http.ListenAndServe(":"+config.PortNumber, api.NewRouter()))

}
