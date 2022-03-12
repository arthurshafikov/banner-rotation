package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig
	ServerConfig
}

type DatabaseConfig struct {
	DSN string
}

type ServerConfig struct {
	Port string
}

func NewConfig(envFileLocation string) *Config {
	err := godotenv.Load(envFileLocation)
	if err != nil {
		panic(err)
	}

	return &Config{
		DatabaseConfig: DatabaseConfig{
			DSN: os.Getenv("DSN"),
		},
		ServerConfig: ServerConfig{
			Port: os.Getenv("APP_PORT"),
		},
	}
}
