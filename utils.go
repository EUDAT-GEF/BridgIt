package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"crypto/tls"
	"io/ioutil"
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

	//if config.StaticContentFolder == "" {
	//	return config, Err(nil, "Empty static content folder name in file: %s", configFilepath)
	//}

	return config, nil
}

// TLSHTTPRequest allows to make requests ignoring the check of certificates
func TLSHTTPRequest(method string, url string, form url.Values) (*http.Response, error) {
	// Ignoring certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := http.Client{Transport: tr}

	req, err := http.NewRequest(method, url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	if form != nil {
		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// StartGEFJob starts a new job in the GEF
func StartGEFJob(serviceID string, accessToken string, pid string) (string, error) {
	log.Println("Starting a new job for the service " + serviceID + " with the PID " + pid)
	// Creating a form
	form := url.Values{}
	form.Add("serviceID", serviceID)
	form.Add("pid", pid)

	resp, err := TLSHTTPRequest("POST", Config.GEFAddress+"/api/jobs?access_token="+accessToken, form)
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
func GetJobStateCode(accessToken string, jobID string) (int, error) {
	resp, err := TLSHTTPRequest("GET", Config.GEFAddress+"/api/jobs/"+jobID+"?access_token="+accessToken, nil)
	if err != nil {
		return 0, err
	}

	// We need to read JSON that normally contains a volumeID
	var jsonReply SelectedJob
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return 0, err
	}

	return jsonReply.Job.State.Code, nil
}

// GetOutputVolumeID returns the output volume ID for the given job
func GetOutputVolumeID(accessToken string, jobID string) (string, error) {
	log.Println("Retrieving the output volume ID for the job " + jobID)
	resp, err := TLSHTTPRequest("GET", Config.GEFAddress+"/api/jobs/"+jobID+"?access_token="+accessToken, nil)
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
func GetVolumeFileName(accessToken string, volumeID string) (string, error) {
	log.Println("Reading the output volume " + volumeID)
	resp, err := TLSHTTPRequest("GET", Config.GEFAddress+"/api/volumes/"+volumeID+"/?access_token="+accessToken, nil)
	if err != nil {
		return "", err
	}
	var jsonReply VolumeInspection

	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return "", err
	}

	if len(jsonReply.VolumeContent) > 0 {
		return filepath.Join(volumeID, jsonReply.VolumeContent[0].Path, jsonReply.VolumeContent[0].Name), nil
	} else {
		return "", nil
	}

}

// GetOutputFile returns a link to the first file (Weblicht service will always produce only one file) from the output volume
func GetOutputFileURL(accessToken string, jobID string) (string, error) {
	log.Println("Retrieving a link to the output file from the job " + jobID)
	for {
		jobExitCode, err := GetJobStateCode(accessToken, jobID)
		if jobExitCode > -1 {
			if err != nil {
				return "", err
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	volumeID, err := GetOutputVolumeID(accessToken, jobID)
	if err != nil {
		return "", err
	}

	fileName, err := GetVolumeFileName(accessToken, volumeID)
	if err != nil {
		return "", err
	}

	return Config.GEFAddress + "/api/volumes/" + fileName + "?content&access_token=" + accessToken, nil
}

// ReadOutputFile reads a file from a certain URL
func ReadOutputFile(fileURL string) ([]byte, error) {
	log.Println("Reading the output file")
	resp, err := TLSHTTPRequest("GET", fileURL, nil)
	if err != nil {
		return nil, err
	}
	outputBuf, err := ioutil.ReadAll(resp.Body)

	return outputBuf, err
}
