package config

import (
	"io/ioutil"
	"path"

	"github.com/go-yaml/yaml"
)

// DefaultConfigFile is the default name of the configuration file.
const DefaultConfigFile = "config.yaml"

// Config represents the blog configuration.
type Config struct {
	Base   string
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

// Load loads the blog configuration.
func Load(folder string) (*Config, error) {
	contents, err := ioutil.ReadFile(path.Join(folder, DefaultConfigFile))
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return nil, err
	}
	cfg.Base = folder
	return &cfg, nil
}
