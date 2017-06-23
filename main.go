package main

import (
	"log"
	"net/http"
)

const configFilePath = "config.json"

func main() {
	config, err := ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}

	router := NewRouter(config)
	log.Fatal(http.ListenAndServe(":"+config.PortNumber, router))
}
