package config

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var apiConfig struct {
		API API
	}

	t.Setenv("API_ADDRESS", "test.com")
	t.Setenv("API_CORS_ORIGINS", "localhost,stage")

	assert.NoError(t, env.Parse(&apiConfig))

	assert.Equal(t, "test.com", apiConfig.API.Address)
	assert.Equal(t, []string{"localhost", "stage"}, apiConfig.API.CORSOrigins)

	var serviceConfig Service

	t.Setenv("SERVICE_NAME", "example")

	assert.NoError(t, env.Parse(&serviceConfig))
	assert.Equal(t, "example", serviceConfig.Name)
}

func TestDB_ConnectionURI(t *testing.T) {
	db := DB{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		Name:     "database_name",
		Schema:   "schema_name",
		SSLMode:  "disable",
	}
	expectedURI := "postgres://user:password@localhost:5432/database_name?search_path=schema_name&sslmode=disable"
	assert.Equal(t, expectedURI, db.ConnectionURI())
}

func TestRabbitMQ_ConnectionURI(t *testing.T) {
	db := RabbitMQ{
		User:        "user",
		Password:    "password",
		Host:        "localhost",
		Port:        "5432",
		VirtualHost: "vhost",
	}
	expectedURI := "amqp://user:password@localhost:5432/vhost"
	assert.Equal(t, expectedURI, db.ConnectionURI())
}

func TestNATS_ConnectionURI(t *testing.T) {
	nats := NATS{
		Host: "localhost",
		Port: "4222",
	}
	assert.Equal(t, "nats://localhost:4222", nats.ConnectionURI())
	nats.User = "user"
	assert.Equal(t, "nats://user@localhost:4222", nats.ConnectionURI())
	nats.Password = "password"
	assert.Equal(t, "nats://user:password@localhost:4222", nats.ConnectionURI())
}
