package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
	Database string `yaml:"database"`
}

type Config struct {
	Port     int      `yaml:"port"`
	Debug    bool     `yaml:"debug"`
	DBConfig DBConfig `yaml:"db"`
}

func NewConfig() *Config {
	data, err := os.ReadFile(getEnv("CONFIG_PATH", "/home/nick/Projects/product-inventory-management-system/cmd/config.yaml"))
	if err != nil {
		// Can not do anything without the config file
		panic(err)
	}

	c := new(Config)
	if err := yaml.Unmarshal(data, c); err != nil {
		panic(err)
	}

	return c
}
