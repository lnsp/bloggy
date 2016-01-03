package config

import (
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	HostAddress string `json:"host"`
	BlogName string `json:"name"`
	HostCountry string `json:"country"`
}

func GetConfig(file string) (Config, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(f, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

