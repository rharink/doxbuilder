package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//Configuration contains typed configuration from file
type Configuration struct {
	OutputDir      string   `yaml:"output_dir"`
	AllowedFormats []string `yaml:"allowed_formats"`
}

//LoadConfiguration loads the configuration from file
func LoadConfiguration(path string) (Configuration, error) {
	c := Configuration{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
