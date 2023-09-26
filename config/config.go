package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type User struct {
	Name        string `yaml:"name"`
	FullName    string `yaml:"full_name"`
	PhoneNumber string `yaml:"phone_number"`
	Email       string `yaml:"email"`
	Duty        []struct {
		Date string `yaml:"date"`
		Role string `yaml:"role"`
	} `yaml:"duty"`
}

type Team struct {
	Name               string `yaml:"name"`
	SchedulingTimezone string `yaml:"scheduling_timezone"`
	Email              string `yaml:"email"`
	SlackChannel       string `yaml:"slack_channel"`
	Users              []User `yaml:"users"`
}

type Config struct {
	Teams []Team `yaml:"teams"`
}

var (
	config Config = Config{}
	once   sync.Once
)

// Get returns the configuration structure initialized with the values ​​of the config.yaml file.
func Get() *Config {
	once.Do(func() {
		configData, err := os.ReadFile("config.yaml")
		if err != nil {
			log.Fatal("Failed to read config file:", err)
		}

		err = yaml.Unmarshal(configData, &config)
		if err != nil {
			log.Fatal("Failed to parse config:", err)
		}
	})
	return &config
}
