package main

import (
	"fmt"

	"io/ioutil"
	"net/http"

	"encoding/json"
	"github.com/pborman/uuid"
	"log"
	"os"
	"path/filepath"
	//"bufio"
)

type Response struct {
	http.ResponseWriter
}

// ClientError sets a 400 error
func (w Response) ClientError(message string, err error) {
	str := fmt.Sprintf("    API Client ERROR: %s\n\t%s", message, err.Error())
	log.Println(str)
	http.Error(w, str, 400)
}

// DirectiveError sets a 403 error
func (w Response) DirectiveError() {
	str := fmt.Sprintf("    API denied by directive ERROR\n")
	log.Println(str)
	http.Error(w, str, 403)
}

// ServerError sets a 500/server error
func (w Response) ServerError(message string, err error) {
	str := fmt.Sprintf("    API Server ERROR: %s\n\t%s", message, err.Error())
	log.Println(str)
	http.Error(w, str, 500)
}

// ServerNewError sets a 500/server error
func (w Response) ServerNewError(message string) {
	str := fmt.Sprintf("    API Server ERROR: %s", message)
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

func Index(w http.ResponseWriter, r *http.Request) {
	Response{w}.Ok(jmap("name", appName, "version", appVersion, "Description", appDescription))
}

func JobStart(w http.ResponseWriter, r *http.Request) {
	var serviceName []string
	var accessToken []string
	var inputFile []string
	var ok bool

	if serviceName, ok = r.URL.Query()["service"]; !ok {
		Response{w}.ServerError("Could not extract a service ID from the request URL", nil)
		return
	}

	if accessToken, ok = r.URL.Query()["token"]; !ok {
		Response{w}.ServerError("Could not extract an access token from the request URL", nil)
		return
	}

	if inputFile, ok = r.URL.Query()["input"]; !ok {
		Response{w}.ServerError("Could not extract an input file name from the request URL", nil)
		return
	}

	content, err := ioutil.ReadFile(inputFile[0])
	uniqueName := uuid.New()

	// Saving the file to serve it to the next Weblicht service
	savedFileName := filepath.Join(Config.StaticContentFolder, uniqueName)
	f, err := os.Create(savedFileName)
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		Response{w}.ServerError("Error while writing in a file", err)
	}

	var serviceID string
	for k := range Config.Apps {
		if k == serviceName[0] {
			serviceID = Config.Apps[k]
		}
	}
	if len(serviceID) == 0 {
		Response{w}.ServerError("Service ID was not found", nil)
	}

	// Making a request to the GEF instance specified in the config file
	jobID, err := StartGEFJob(serviceID, accessToken[0], Config.StorageURL+":"+Config.StoragePortNumber+Config.StaticContentURLPrefix+"/"+uniqueName)

	if err != nil {
		Response{w}.ServerError("Error while starting a new job", err)
	}

	outputFileLink, err := GetOutputFileURL(accessToken[0], jobID)
	if err != nil {
		Response{w}.ServerError("Error while getting a link to the output file", err)
	}

	outputBuf, err := ReadOutputFile(outputFileLink)
	if err != nil {
		Response{w}.ServerError("Error while reading the output file", err)
	}
	outputType := http.DetectContentType(outputBuf)

	Response{w}.Header().Set("Content-Type", outputType+"+url")
	Response{w}.Write([]byte(outputFileLink))
	//w.Header().Set("Content-Type", outputType+"+url")
	//w.Write([]byte(outputFileLink))
}
