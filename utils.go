package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"crypto/tls"
	"log"
	"path/filepath"
)

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
	log.Println("Starting a new job for the service " + serviceID + " with the PID " + pid)
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

	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}

	// We need to read JSON that normally contains a jobID
	var jsonReply map[string]interface{}
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

// GetJobStateCode returns the job exit code (-1 running, 0 ended successfully, 1 failed)
func GetJobStateCode(jobID string) (int, error) {
	//log.Println("Reading the state of the job " + jobID)

	// Ignoring certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := http.Client{Transport: tr}

	routerURL := Config.GEFAddress + "/api/jobs/" + jobID // GEF endpoint

	// Sending a GET request
	req, err := http.NewRequest("GET", routerURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := hc.Do(req)
	if err != nil {
		return 0, err
	}
	var jsonReply SelectedJob

	// We need to read JSON that normally contains a volumeID
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return 0, err
	}

	return jsonReply.Job.State.Code, nil
}

// GetOutputVolumeID returns the output volume ID for the given job
func GetOutputVolumeID(jobID string) (string, error) {
	log.Println("Retrieving the output volume ID for the job " + jobID)
	// Ignoring certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := http.Client{Transport: tr}

	routerURL := Config.GEFAddress + "/api/jobs/" + jobID // GEF endpoint

	// Sending a GET request
	req, err := http.NewRequest("GET", routerURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}

	var jsonReply SelectedJob
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return "", err
	}

	return jsonReply.Job.OutputVolume, nil
}

// GetVolumeFile inspects the output volume and return a path to the output file
func GetVolumeFile(volumeID string) (string, error) {
	log.Println("Reading the output volume " + volumeID)
	// Ignoring certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := http.Client{Transport: tr}

	routerURL := Config.GEFAddress + "/api/volumes/" + volumeID + "/" // GEF endpoint

	// Sending a GET request
	req, err := http.NewRequest("GET", routerURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}

	var jsonReply VolumeInspection

	err = json.NewDecoder(resp.Body).Decode(&jsonReply)

	if len(jsonReply.VolumeContent) > 0 {
		return filepath.Join(volumeID, jsonReply.VolumeContent[0].Path, jsonReply.VolumeContent[0].Name), nil
	}

	return "", err
}

// GetOutputFile returns a link to the first file (Weblicht service will always produce only one file) from the output volume
func GetOutputFile(jobID string) (string, error) {
	log.Println("Retrieving a link to the output file from the job " + jobID)
	for {
		jobExitCode, err := GetJobStateCode(jobID)
		//fmt.Println(jobExitCode)
		if jobExitCode > -1 {
			if err != nil {
				return "", err
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	volumeID, err := GetOutputVolumeID(jobID)
	if err != nil {
		return "", err
	}

	fileName, err := GetVolumeFile(volumeID)
	if err != nil {
		return "", err
	}

	return Config.GEFAddress + "/api/volumes/" + fileName + "?content", nil
}
