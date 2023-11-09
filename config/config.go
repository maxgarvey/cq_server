package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Redis    Redis    `yaml:"redis"`
	Rabbitmq Rabbitmq `yaml:"rabbitmq"`
	Server   Server   `yaml:"server"`
}

type Rabbitmq struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Queuename string `yaml:"queuename"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	Port int `yaml:"port"`
}

func GetConfig(configFile string) Config {
	var c Config

	// Read file to byte array.
	yamlFile, err := os.ReadFile(
		configFile,
	)
	if err != nil {
		log.Printf(
			"Error reading config: %v ",
			err,
		)
	}

	// Unmarshal into config struct.
	err = yaml.Unmarshal(
		yamlFile,
		&c,
	)
	if err != nil {
		log.Fatalf(
			"Error unmarshalling YAML: %v",
			err,
		)
	}

	return c
}
