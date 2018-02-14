package utils

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

	"github.com/EUDAT-GEF/BridgIt/def"
)

// ReadConfigFile reads a configuration file
func ReadConfigFile(configFilepath string) (def.Configuration, error) {
	var config def.Configuration

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
func StartGEFJob(serviceID string, accessToken string, pid string, GEFAddress string) (string, error) {
	log.Println("Starting a new job for the service " + serviceID + " with the PID " + pid)

	// Creating a form
	form := url.Values{}
	form.Add("serviceID", serviceID)
	form.Add("pid", pid)

	resp, err := TLSHTTPRequest("POST", GEFAddress+"/api/jobs?access_token="+accessToken, form)
	if err != nil {
		return "", Err(err, "Failed to send a POST request that starts a job")
	}

	// We need to read JSON that normally contains a jobID
	var jsonReply map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return "", Err(err, "Failed to parse the GEF server reply")
	}

	if val, ok := jsonReply["jobID"]; ok {
		if jobID, ok := val.(string); ok {
			return jobID, nil
		} else {
			return "", Err(nil, "Failed to convert the output to string")
		}
	} else {
		return "", Err(nil, "JobID was not found in the GEF server reply")
	}
}

// GetJobStateCode returns the job exit code (-1 running, 0 ended successfully, 1 failed)
func GetJobStateCode(accessToken string, jobID string, GEFAddress string) (int, error) {
	resp, err := TLSHTTPRequest("GET", GEFAddress+"/api/jobs/"+jobID+"?access_token="+accessToken, nil)

	if err != nil {
		return 1, Err(err, "Failed to send a GET request that returns a job status")
	}

	var jsonReply def.SelectedJob
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	log.Print("REPLY")
	log.Print(resp.Body)
	if err != nil {
		return 1, Err(err, "Failed to parse the GEF server reply")
	}

	return jsonReply.Job.State.Code, nil
}

// GetOutputVolumeID returns the output volume ID for the given job
func GetOutputVolumeID(accessToken string, jobID string, GEFAddress string) (string, error) {
	log.Println("Retrieving the output volume ID for the job " + jobID)
	resp, err := TLSHTTPRequest("GET", GEFAddress+"/api/jobs/"+jobID+"?access_token="+accessToken, nil)

	if err != nil {
		return "", Err(err, "Failed to get an output volume id for the job %s", jobID)
	}

	var jsonReply def.SelectedJob
	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	log.Print("Reply")
	log.Print(resp.Body)
	if err != nil {
		return "", Err(err, "Failed to parse the GEF server reply")
	}

	return jsonReply.Job.OutputVolume[0].VolumeID, nil
}

// GetVolumeFileName inspects the output volume and return a path to the output file
func GetVolumeFileName(accessToken string, volumeID string, GEFAddress string) (string, error) {
	log.Println("Reading the output volume " + volumeID)
	resp, err := TLSHTTPRequest("GET", GEFAddress+"/api/volumes/"+volumeID+"/?access_token="+accessToken, nil)
	if err != nil {
		return "", Err(err, "Failed to get a path to the output file in the volume %s", volumeID)
	}
	var jsonReply def.VolumeInspection

	err = json.NewDecoder(resp.Body).Decode(&jsonReply)
	if err != nil {
		return "", Err(err, "Failed to parse the GEF server reply")
	}

	if len(jsonReply.VolumeContent) > 0 {
		return filepath.Join(volumeID, jsonReply.VolumeContent[0].Path, jsonReply.VolumeContent[0].Name), nil
	} else {
		return "", nil
	}
}

// GetOutputFileURL returns a link to the first file (Weblicht service will always produce only one file) from the output volume
func GetOutputFileURL(accessToken string, jobID string, GEFAddress string) (string, error) {
	log.Println("Retrieving a link to the output file from the job " + jobID)
	for {
		jobExitCode, err := GetJobStateCode(accessToken, jobID, GEFAddress)
		if jobExitCode > -1 {
			if err != nil {
				return "", Err(err, "The job failed")
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Print("GEtting output volume")
	volumeID, err := GetOutputVolumeID(accessToken, jobID, GEFAddress)
	if err != nil {
		return "", err
	}

	fileName, err := GetVolumeFileName(accessToken, volumeID, GEFAddress)
	if err != nil {
		return "", err
	}
	log.Print("Reply")
	log.Print(GEFAddress + "/api/volumes/" + fileName + "?content&access_token=" + accessToken)
	return GEFAddress + "/api/volumes/" + fileName + "?content&access_token=" + accessToken, nil
}

// ReadOutputFile reads a file from a certain URL
func ReadOutputFile(fileURL string) ([]byte, error) {
	log.Println("Reading the output file: " + fileURL)
	resp, err := TLSHTTPRequest("GET", fileURL, nil)
	if err != nil {
		return nil, Err(err, "Could not access the output file %s", fileURL)
	}
	outputBuf, err := ioutil.ReadAll(resp.Body)

	return outputBuf, err
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
