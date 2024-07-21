package cmd

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Algorithm string `yaml:"algorithm"`
	Port      int    `yaml:"port"`
	Strict    bool   `yaml:"strict"`
	Servers   []struct {
		Host        string  `yaml:"host"`
		Weight      float64 `yaml:"weight"`
		Connections int     `yaml:"connections"`
	} `yaml:"servers"`
}

func (c *Config) GetConf() *Config {
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
