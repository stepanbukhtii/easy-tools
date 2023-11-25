package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testConfig struct {
	API API `mapstructure:"API"`
}

func TestInit(t *testing.T) {
	var config testConfig
	assert.NoError(t, os.Setenv("API_URL", "test.com"))
	assert.NoError(t, Load(&config))
	fmt.Println("config", config)
}
