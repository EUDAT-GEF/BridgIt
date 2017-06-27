package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"crypto/tls"
	"reflect"
)

func PseudoUUID() string {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	unixTime := int64(time.Now().Unix())
	uniqueName := uuid + "-" + strconv.FormatInt(unixTime, 10)

	return uniqueName
}

// ReadConfigFile reads a configuration file
func ReadConfigFile(configFilepath string) (Configuration, error) {
	var config Configuration

	file, err := os.Open(configFilepath)
	if err != nil {
		return config, Err(err, "Cannot open config file %s", configFilepath)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, Err(err, "Cannot read config file %s", configFilepath)
	}

	if config.StaticContent == "" {
		return config, Err(nil, "Empty static content folder name in file: %s", configFilepath)
	}

	return config, nil
}

// StartGEFJob starts a new job in the GEF
func StartGEFJob(serviceID string, pid string) (string, error) {
	// Ignoring certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := http.Client{Transport: tr}

	routerURL := Config.GEFAddress + "/api/jobs" // GEF endpoint
	// Creating a form
	form := url.Values{}
	form.Add("serviceID", serviceID)
	form.Add("pid", pid)

	// POSTing the request
	req, err := http.NewRequest("POST", routerURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Processing the reply
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	var jsonReply map[string]interface{}
	// We need to read JSON that normally contains a jobID
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return "", err
	}

	if val, ok := jsonReply["jobID"]; ok {
		if jobID, ok := val.(string); ok {
			return jobID, nil
		}
	}

	return "", Err(err, "Failed to convert the output to string")
}


type Job struct {
	ID           string
	ServiceID    string
	Input        string
	Created      time.Time
	State        *JobState
	InputVolume  string
	OutputVolume string
	Tasks        []Task
}

// JobState keeps information about a job state
type JobState struct {
	Status string
	Error  string
	Code   int
}

// Task contains tasks related to a specific job (used to serialize JSON)
type Task struct {
	ID            string
	Name          string
	ContainerID   string
	Error         string
	ExitCode      int
	ConsoleOutput string
}

func Keys(v interface{}) ([]string, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return nil, Err(nil, "not a map")
	}
	t := rv.Type()
	if t.Key().Kind() != reflect.String {
		return nil, Err(nil, "not string key")
	}
	var result []string
	for _, kv := range rv.MapKeys() {
		result = append(result, kv.String())
	}
	return result, nil
}

func GetJobStateCode(job interface{}) (float64, error) {
	jobMap := reflect.ValueOf(job)

	if jobMap.Kind() == reflect.Map {
		for _, jobKey := range jobMap.MapKeys() {
			jobItem := jobMap.MapIndex(jobKey)
			if strings.ToLower(jobKey.String()) == "state" {
				stateMap := reflect.ValueOf(jobItem.Interface())

				if stateMap.Kind() == reflect.Map {
					for _, stateKey := range stateMap.MapKeys() {
						state := stateMap.MapIndex(stateKey)
						if strings.ToLower(stateKey.String()) == "code" {
							if val, ok := state.Interface().(float64); ok {
								return val, nil
							} else {
								return 0, Err(nil, "Failed to convert the output to float64")
							}
						}
					}
				}
			}
		}
	}

	return 0, nil
}

func GetOutputVolumeID(jobID string) (string, error) {
	// Ignoring certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := http.Client{Transport: tr}

	routerURL := Config.GEFAddress + "/api/jobs/"+ jobID // GEF endpoint

	// Sending a GET request
	req, err := http.NewRequest("GET", routerURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := hc.Do(req)


	// Processing the reply

	if err != nil {
		return "", err
	}
	var jsonReply map[string]interface{}

	// We need to read JSON that normally contains a volumeID
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return "", err
	}
	fmt.Println(jsonReply)
	fmt.Println("JOB = ")
	if job, ok := jsonReply["Job"]; ok {
		
		fmt.Println(job)


		fmt.Println("CODE = ")
		fmt.Println(GetJobStateCode(job))

	}


	return "", Err(err, "Failed to convert the output to string")
}