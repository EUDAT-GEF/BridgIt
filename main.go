package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
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

func startHttpServer(portNumber string, router *mux.Router) *http.Server {
	log.Println("Starting the service at port " + portNumber)
	srv := &http.Server{Addr: ":"+portNumber, Handler: router}

	go func() {
		log.Println(http.ListenAndServe(":"+portNumber, router))
		//if err := srv.ListenAndServe(); err != nil {
		//	log.Printf("Httpserver: ListenAndServe() error: %s", err)
		//}
	}()

	return srv
}

func main() {
	Config = initApp()
	//startHttpServer(Config.PortNumber, NewRouter())
	//
	log.Fatal(http.ListenAndServe(":"+Config.PortNumber, NewRouter()))

}
