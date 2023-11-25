package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Service struct {
	Name        string `env:"SERVICE_NAME"`
	Environment string `env:"SERVICE_ENVIRONMENT"`
	Version     string `env:"SERVICE_VERSION"`
}

type API struct {
	Address        string        `env:"API_ADDRESS"`
	CORSOrigins    []string      `env:"API_CORS_ORIGINS"`
	JWT            JWT           `env:"API_JWT"`
	Timeout        time.Duration `env:"API_TIMEOUT"`
	SwaggerEnabled bool          `env:"API_SWAGGER_ENABLED"`
}

type JWT struct {
	Enabled    bool          `env:"JWT_ENABLED"`
	PublicKey  string        `env:"JWT_PUBLIC_KEY"`
	PrivateKey string        `env:"JWT_PRIVATE_KEY"`
	ClaimsTTL  time.Duration `env:"JWT_CLAIMS_TTL"`
}

type Log struct {
	Level string `env:"LOG_LEVEL"`
}

func (l Log) IsDebug() bool {
	return strings.EqualFold(l.Level, zerolog.LevelDebugValue)
}

type DB struct {
	Host               string `env:"DB_HOST"`
	Port               string `env:"DB_PORT"`
	User               string `env:"DB_USER"`
	Password           string `env:"DB_PASSWORD"`
	Name               string `env:"DB_NAME"`
	Schema             string `env:"DB_SCHEMA"`
	SSLMode            string `env:"DB_SSL_MODE"`
	MaxOpenConnections *int   `env:"DB_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections *int   `env:"DB_MAX_IDLE_CONNECTIONS"`
}

func (db DB) ConnectionString() string {
	s := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.Name)

	var params []string
	if db.Schema != "" {
		params = append(params, fmt.Sprintf("search_path=%s", db.Schema))
	}
	if db.SSLMode != "" {
		params = append(params, fmt.Sprintf("sslmode=%s", db.SSLMode))
	}
	if len(params) > 0 {
		s = fmt.Sprintf("%s?%s", s, strings.Join(params, "&"))
	}

	return s
}

type Redis struct {
	Addresses   []string `env:"REDIS_ADDRESSES"`
	MasterName  string   `env:"REDIS_MASTER"`
	Password    string   `env:"REDIS_PASSWORD"`
	DB          int      `env:"REDIS_DB"`
	TLSDisabled bool     `env:"REDIS_TLS_DISABLED"`
}

type NATS struct {
	Host     string `env:"NATS_HOST"`
	Port     string `env:"NATS_PORT"`
	User     string `env:"NATS_USER"`
	Password string `env:"NATS_PASSWORD"`
}
