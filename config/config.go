package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Postgres Postgres `yaml:"postgres" json:"postgres"`
	Redis    Redis    `yaml:"redis" json:"redis"`
	Rabbitmq Rabbitmq `yaml:"rabbitmq" json:"rabbitmq"`
	Server   Server   `yaml:"server" json:"server"`
}

type Postgres struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DBName   string `yaml:"db_name" json:"db_name"`
}

type Rabbitmq struct {
	Username  string `yaml:"username" json:"username"`
	Password  string `yaml:"password" json:"password"`
	Host      string `yaml:"host" json:"host"`
	Port      int    `yaml:"port" json:"port"`
	Queuename string `yaml:"queuename" json:"queuename"`
}

type Redis struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
}

type Server struct {
	Port int `yaml:"port" json:"port"`
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
