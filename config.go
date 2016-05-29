package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

// DefaultConfigFile is the default name of the configuration file.
const DefaultConfigFile = "config.json"

// BlogConfig represents the blog configuration.
type BlogConfig struct {
	Address    string            `json:"host"`
	AddressTLS string            `json:"host.secure"`
	Country    string            `json:"country"`
	Title      string            `json:"title"`
	Subtitle   string            `json:"subtitle"`
	Author     string            `json:"author"`
	Email      string            `json:"email"`
	URL        string            `json:"url"`
	Links      map[string]string `json:"links"`
}

// LoadConfig loads the blog configuration.
func LoadConfig() error {
	contents, readError := ioutil.ReadFile(path.Join(BlogFolder, DefaultConfigFile))
	if readError != nil {
		return readError
	}

	jsonError := json.Unmarshal(contents, &GlobalConfig)
	if jsonError != nil {
		return jsonError
	}

	return nil
}
