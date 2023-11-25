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
