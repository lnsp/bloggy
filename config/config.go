package config

import (
	"io/ioutil"
	"encoding/json"
)

type Configuration struct {
	HostPort int `json:"port"`
	BlogName string `json:"name"`
	HostCountry string `json:"country"`
}

func GetConfig(file string) (Configuration, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return Configuration{}, err
	}

	var config Configuration
	err = json.Unmarshal(f, &config)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}

func GetDefaultConfig() (Configuration, error) {
	return GetConfig("./config.json")
}

