package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig
	ServerConfig
	MultihandedBanditConfig
	QueueConfig
}

type DatabaseConfig struct {
	DSN string
}

type ServerConfig struct {
	Port string
}

type MultihandedBanditConfig struct {
	EGreedValue float64
}

type QueueConfig struct {
	BrokerAddress string
}

func NewConfig(envFileLocation string) *Config {
	err := godotenv.Load(envFileLocation)
	if err != nil {
		panic(err)
	}

	eGreedValue, err := strconv.ParseFloat(os.Getenv("E_GREED_VALUE"), 64)
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
		MultihandedBanditConfig: MultihandedBanditConfig{
			EGreedValue: eGreedValue,
		},
		QueueConfig: QueueConfig{
			BrokerAddress: os.Getenv("BROKER_ADDRESS"),
		},
	}
}
