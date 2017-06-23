package main

import (
	"fmt"
	"time"
	"crypto/rand"
	"os"
	"encoding/json"
	"strconv"
)

func PseudoUUID() (string) {

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