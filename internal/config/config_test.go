package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	envFileContent = []byte(
		`POSTGRES_DB=homestead
		POSTGRES_USER=homestead
		POSTGRES_PASSWORD=secret
		APP_PORT=8123
		DSN=host=db user=homestead password=secret dbname=homestead sslmode=disable
		E_GREED_VALUE=0.1
		BROKER_ADDRESS=kafka:9092`,
	)
	envFilePath = ".env.test"

	expectedDatabaseConfig = DatabaseConfig{
		DSN: "host=db user=homestead password=secret dbname=homestead sslmode=disable",
	}
	expectedServerConfig = ServerConfig{
		Port: "8123",
	}
	expectedMultihandedBanditConfig = MultihandedBanditConfig{
		EGreedValue: 0.1,
	}
	expectedQueueConfig = QueueConfig{
		BrokerAddress: "kafka:9092",
	}
)

func TestNewConfigFromEnvFile(t *testing.T) {
	createFakeEnvFile(t)
	defer deleteFakeEnvFile(t)

	config := NewConfig(envFilePath)

	require.Equal(t, expectedDatabaseConfig, config.DatabaseConfig)
	require.Equal(t, expectedServerConfig, config.ServerConfig)
	require.Equal(t, expectedMultihandedBanditConfig, config.MultihandedBanditConfig)
	require.Equal(t, expectedQueueConfig, config.QueueConfig)
}

func createFakeEnvFile(t *testing.T) {
	t.Helper()
	if err := os.WriteFile(envFilePath, envFileContent, 0600); err != nil { //nolint:gofumpt
		t.Fatal(err)
	}
}

func deleteFakeEnvFile(t *testing.T) {
	t.Helper()
	if err := os.Remove(envFilePath); err != nil {
		t.Fatal(err)
	}
}
