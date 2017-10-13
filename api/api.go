package api

import (
	"fmt"

	"net/http"

	"encoding/json"
	"log"

	"github.com/EUDAT-GEF/Bridgit/def"
	"github.com/EUDAT-GEF/Bridgit/utils"

	"github.com/gorilla/mux"
)

type Response struct {
	http.ResponseWriter
}

type App struct {
	Info   def.Info
	Config def.Configuration
	Server *http.Server
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// ServerError sets a 500/server error
func (w Response) ServerError(message string, err error) {
	str := fmt.Sprintf("API Server error: %s\n\t%s", message, err.Error())
	log.Println(str)
	http.Error(w, str, 500)
}

// ServerNewError sets a 500/server error
func (w Response) ServerNewError(message string) {
	str := fmt.Sprintf("API Server error: %s", message)
	log.Println(str)
	http.Error(w, str, 500)
}

// Location sets location header
func (w Response) Location(loc string) Response {
	w.Header().Set("Location", loc)
	return w
}

// Ok sets 200/ok response code and body
func (w Response) Ok(body interface{}) {
	setCodeAndBody(w, 200, body)
}

// Created sets 201/created response code and body
func (w Response) Created(body interface{}) {
	setCodeAndBody(w, 201, body)
}

func setCodeAndBody(w Response, code int, body interface{}) {
	var contentType string
	var data []byte
	var err error

	if jsonMap, ok := body.(map[string]interface{}); ok {
		data, err = json.Marshal(jsonMap)
		if err != nil {
			w.ServerError("json marshal: ", err)
			http.Error(w, fmt.Sprintln("Server Error: unexpected Ok body type"), 500)
			return
		}
		contentType = "application/json; charset=utf-8"
	} else if str, ok := body.(string); ok {
		contentType = "text/plain; charset=utf-8"
		data = []byte(str)
	} else {
		log.Printf("ERROR: unexpected Ok body type: %T\n", body)
		http.Error(w, fmt.Sprintln("Server Error: unexpected Ok body type"), 500)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	w.Write(data)
	// log.Println("setCodeAndBody:", code, contentType, body)
	log.Println(" -> HTTP", code, ",", contentType, ",", len(data), "bytes")
}

func jmap(kv ...interface{}) map[string]interface{} {
	if len(kv) == 0 {
		log.Println("ERROR: jsonmap: empty call")
		return nil
	} else if len(kv)%2 == 1 {
		log.Println("ERROR: jsonmap: unbalanced call")
		return nil
	}
	m := make(map[string]interface{})
	k := ""
	for _, kv := range kv {
		if k == "" {
			if skv, ok := kv.(string); !ok {
				log.Println("ERROR: jsonmap: expected string key")
				return m
			} else if skv == "" {
				log.Println("ERROR: jsonmap: string key is empty")
				return m
			} else {
				k = skv
			}
		} else {
			m[k] = kv
			k = ""
		}
	}
	return m
}

// NewApp creates a Bridgit application object (with config and server information) to initialize a service
func NewApp(cfg def.Configuration) App {
	log.Println("Preparing to start the service at port " + cfg.PortNumber)
	srv := &http.Server{Addr: ":" + cfg.PortNumber}

	application := App{
		Info: def.Info{
			Name:        "BridgIt",
			Version:     "0.1",
			Description: "This is Bridgit, a liaison between Weblicht and the GEF",
		},
		Config: cfg,
		Server: srv,
	}

	var routes = Routes{
		Route{
			"Index",
			"GET",
			"/",
			application.Index,
		},
		Route{
			"JobStart",
			"POST",
			"/jobs",
			application.JobStart,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = utils.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	application.Server.Handler = router
	return application
}

// Start starts Bridgit service
func (a *App) Start() error {
	log.Println("Starting the service...")
	return a.Server.ListenAndServe()
}

// Stop stops Bridgit service
func (a *App) Stop() error {
	return a.Server.ListenAndServe()
}

// Index shows information about the API
func (a *App) Index(w http.ResponseWriter, r *http.Request) {
	Response{w}.Ok(jmap("name", a.Info.Name, "version", a.Info.Version, "Description", a.Info.Description))
}

// JobStart starts a job at the GEF instance
func (a *App) JobStart(w http.ResponseWriter, r *http.Request) {
	var serviceName []string
	var accessToken []string
	var inputFile []string
	var ok bool

	if serviceName, ok = r.URL.Query()["service"]; !ok {
		Response{w}.ServerNewError("Could not extract a service name from the request URL")
		return
	}

	if accessToken, ok = r.URL.Query()["token"]; !ok {
		Response{w}.ServerNewError("Could not extract an access token from the request URL")
		return
	}

	if inputFile, ok = r.URL.Query()["input"]; !ok {
		Response{w}.ServerNewError("Could not extract an input file name from the request URL")
		return
	}

	var serviceID string
	serviceFound := false
	for k := range a.Config.Apps {
		if k == serviceName[0] {
			serviceFound = true
			serviceID = a.Config.Apps[k]
			break
		}
	}

	if !serviceFound {
		Response{w}.ServerNewError("Could not locate any services corresponding to the service name specified in the request URL")
		return
	}

	// Making a request to the GEF instance specified in the config file
	jobID, err := utils.StartGEFJob(serviceID, accessToken[0], inputFile[0], a.Config.GEFAddress)
	if err != nil {
		Response{w}.ServerError("Error while starting a new job on the GEF instance", err)
		return
	}

	outputFileLink, err := utils.GetOutputFileURL(accessToken[0], jobID, a.Config.GEFAddress)
	if err != nil {
		Response{w}.ServerError("Error while getting a link to the output file", err)
		return
	}

	outputBuf, err := utils.ReadOutputFile(outputFileLink)
	if err != nil {
		Response{w}.ServerError("Error while reading the output file", err)
		return
	}
	outputType := http.DetectContentType(outputBuf)

	Response{w}.Header().Set("Content-Type", outputType+"+url")
	Response{w}.Write([]byte(outputFileLink))
}
