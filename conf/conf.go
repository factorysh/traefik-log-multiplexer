package conf

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Path    string
	Input   map[string]map[string]interface{}   `yaml:"input"`
	Filters []map[string]map[string]interface{} `yaml:"filters"`
	Output  map[string]map[string]interface{}   `yaml:"output"`
}

func Read() (*Config, error) {
	// read yaml config file
	// get path from env
	path := os.Getenv("CONFIG")
	if path == "" {
		path = "/etc/traefik-demultiplexer.yml"
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Error: can't open config file : " + path)
	}

	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, errors.New("Error: cant read config file : " + path)
	}
	c.Path = path

	return &c, nil
}
