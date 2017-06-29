package main

import (
	"log"
	"net/http"
)

const configFilePath = "config.json"
const appName = "BridgIt"
const appVersion = "0.1"
const appDescription = "This is Bridgit, a liaison between Weblicht and the GEF"

var Config Configuration

func initApp() Configuration {
	Config, err := ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	return Config
}

func main() {
	Config = initApp()

	router := NewRouter()
	log.Println("Starting the service at port " + Config.PortNumber)
	log.Fatal(http.ListenAndServe(":"+Config.PortNumber, router))
}
