package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	InNamespace struct {
		Service string `yaml:"service"`
	} `yaml:"in-namespace"`
	CrossNamespace struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Service  string `yaml:"service"`
		DbName   string `yaml:"db-name"`
	} `yaml:"cross-namespace"`
}

func ReadFromYaml(f string) (c *Config, err error) {
	file, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
