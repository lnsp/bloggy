package main

import (
	"io/ioutil"
	"path"

	"github.com/go-yaml/yaml"
)

// DefaultConfigFile is the default name of the configuration file.
const DefaultConfigFile = "config.yaml"

// BlogConfig represents the blog configuration.
type BlogConfig struct {
	Server struct {
		Port    int
		TLSPort int
	}
	Meta struct {
		Country  string
		Title    string
		Subtitle string
		Favicon  string
	}
	Author struct {
		Name  string
		Email string
	}
	Links map[string]string
}

// LoadConfig loads the blog configuration.
func LoadConfig() error {
	contents, err := ioutil.ReadFile(path.Join(BlogFolder, DefaultConfigFile))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(contents, &Config)
	if err != nil {
		return err
	}

	return nil
}
