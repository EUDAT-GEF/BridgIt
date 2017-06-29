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
		if jobID, ok := val.(string); ok {
			return jobID, nil
		}
	}

	return "", Err(err, "Failed to convert the output to string")
}

type SelectedJob struct {
	Job SingleJob
}
type SingleJob struct {
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


func GetJobStateCode(jobID string) (int, error) {
	// Ignoring certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := http.Client{Transport: tr}

	routerURL := Config.GEFAddress + "/api/jobs/"+ jobID // GEF endpoint

	// Sending a GET request
	req, err := http.NewRequest("GET", routerURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := hc.Do(req)
	if err != nil {
		return 0, err
	}
	//var jsonReply map[string]interface{}
	var jsonReply SelectedJob

	// We need to read JSON that normally contains a volumeID
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return 0, err
	}

	fmt.Println(jsonReply.Job.State.Code)
	return jsonReply.Job.State.Code, nil
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
	if err != nil {
		return "", err
	}
	//var jsonReply map[string]interface{}

	var jsonReply SelectedJob
	// We need to read JSON that normally contains a volumeID
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return "", err
	}


	fmt.Println(jsonReply.Job.OutputVolume)
	return jsonReply.Job.OutputVolume, nil
}
type VolumeInspection struct {
	VolumeContent []VolumeItem
}

type VolumeItem struct {
	Name       string       `json:"name"`
	Size       int64        `json:"size"`
	Modified   time.Time    `json:"modified"`
	IsFolder   bool         `json:"isFolder"`
	Path   	   string       `json:"path"`
	FolderTree []VolumeItem `json:"folderTree"`
}

func GetVolumeFile(volumeID string) (string, error) {
	// Ignoring certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := http.Client{Transport: tr}

	routerURL := Config.GEFAddress + "/api/volumes/"+ volumeID + "/" // GEF endpoint

	// Sending a GET request
	req, err := http.NewRequest("GET", routerURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	//var jsonReply map[string]interface{}
	var jsonReply VolumeInspection


	//jsonReceived, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//
	//}


	err = json.NewDecoder(resp.Body).Decode(&jsonReply)




	fmt.Println(jsonReply)
	fmt.Println(jsonReply.VolumeContent)

	if len(jsonReply.VolumeContent)>0 {
		fmt.Println(jsonReply.VolumeContent[0].Name)
	}
	fmt.Println(err)














	return "", nil

}


func GetOutputFile(jobID string) (string, error) {


	for {
		jobExitCode, err := GetJobStateCode(jobID)
		fmt.Println(jobExitCode)
		if jobExitCode >-1 {
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
	fmt.Println("VolumeID = ")
	fmt.Println(volumeID)

	fileName, err := GetVolumeFile(volumeID)
	fmt.Println("Volume File = ")
	fmt.Println(fileName)
	fmt.Println(err)

	return "", nil
}