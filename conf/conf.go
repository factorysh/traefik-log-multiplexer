package conf

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Path    string
	Input   map[string]map[string]interface{}   `yaml:"input"`
	Filters []map[string]map[string]interface{} `yaml:"filters"`
	Output  map[string]map[string]interface{}   `yaml:"output"`
	Admin   *Admin                              `yaml:"admin"`
}

type Admin struct {
	Listen     string `yaml:"listen"`
	Prometheus bool   `yaml:"prometheus"`
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

	c := Config{
		Admin: &Admin{},
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, errors.New("Error: cant read config file : " + path)
	}
	c.Path = path

	err = mergo.Merge(c.Admin, Admin{
		Listen:     "localhost:8345",
		Prometheus: true,
	})
	if err != nil {
		return nil, err
	}

	return &c, nil
}
