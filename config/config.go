package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Redis struct {
		ConnectionType string `yaml:"connection_type"`
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
	} `yaml:"redis"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func GetConfig(environment string) Config {
	var c Config

	// Read file to byte array.
	yamlFile, err := os.ReadFile(
		fmt.Sprintf("%s.yaml", environment))
	if err != nil {
		log.Printf("Error reading config: %v ", err)
	}

	// Unmarshal into config struct.
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	return c
}
