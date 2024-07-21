package cmd

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Algorithm string `yaml:"algorithm"`
	Port      int    `yaml:"port"`
	Servers   []struct {
		Host        string  `yaml:"host"`
		Weight      float64 `yaml:"weight"`
		Connections int     `yaml:"connections"`
	} `yaml:"servers"`
}

func (c *Config) GetConf() *Config {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func SaveLogs() {
	f, err := os.OpenFile("logs/tb.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}
