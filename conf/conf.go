package conf

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogPath string
	Output  map[string]map[string]interface{}
}

func New() (*Config, error) {
	// read yaml config file
	// get path from env
	path := os.Getenv("CONFIG")
	if path == "" {
		path = "/etc/traefik-demultiplexer.yml"
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Error: can't open gitlab config file : " + path)
	}

	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, errors.New("Error: cant read gitlab config file : " + path)
	}

	return &c, nil
}
