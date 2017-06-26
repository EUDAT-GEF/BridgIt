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
		fmt.Println(val)
		if jobID, ok := val.(string); ok {
			return jobID, nil
		}
	}

	return "", Err(err, "Failed to convert the output to string")
}
