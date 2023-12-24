package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	var apiConfig struct {
		API API `mapstructure:"API"`
	}
	assert.NoError(t, os.Setenv("API_ADDRESS", "test.com"))
	assert.NoError(t, os.Setenv("API_CORS_ORIGINS", "localhost,stage"))
	assert.NoError(t, Load(&apiConfig))
	assert.Equal(t, "test.com", apiConfig.API.Address)
	assert.Equal(t, []string{"localhost", "stage"}, apiConfig.API.CORSOrigins)

	var appConfig Application
	assert.NoError(t, os.Setenv("SERVICE_NAME", "example"))
	assert.NoError(t, Load(&appConfig))
	assert.Equal(t, "example", appConfig.ServiceName)
}
