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
	HostAddress    string `json:"host"`
	HostAddressTLS string `json:"host.secure"`
	HostCountry    string `json:"country"`
	BlogTitle      string `json:"title"`
	BlogSubtitle   string `json:"subtitle"`
	BlogAuthor     string `json:"author"`
	BlogEmail      string `json:"email"`
	BlogURL        string `json:"url"`
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
