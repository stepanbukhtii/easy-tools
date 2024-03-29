package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Application struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
	Version     string `mapstructure:"VERSION"`
}

type API struct {
	Address        string        `mapstructure:"ADDRESS"`
	CORSOrigins    []string      `mapstructure:"CORS_ORIGINS"`
	JWT            JWT           `mapstructure:"JWT"`
	Timeout        time.Duration `mapstructure:"TIMEOUT"`
	SwaggerEnabled bool          `mapstructure:"SWAGGER_ENABLED"`
}

type JWT struct {
	Enabled    bool          `mapstructure:"ENABLED"`
	PublicKey  string        `mapstructure:"PUBLIC_KEY"`
	PrivateKey string        `mapstructure:"PRIVATE_KEY"`
	ClaimsTTL  time.Duration `mapstructure:"CLAIMS_TTL"`
}

type Log struct {
	Level string `mapstructure:"LEVEL"`
}

func (l Log) IsDebug() bool {
	return strings.EqualFold(l.Level, zerolog.LevelDebugValue)
}

type DB struct {
	Host               string `mapstructure:"HOST"`
	Port               string `mapstructure:"PORT"`
	User               string `mapstructure:"USER"`
	Password           string `mapstructure:"PASSWORD"`
	Name               string `mapstructure:"NAME"`
	Schema             string `mapstructure:"SCHEMA"`
	SSLMode            string `mapstructure:"SSL_MODE"`
	MaxOpenConnections *int   `mapstructure:"MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections *int   `mapstructure:"MAX_IDLE_CONNECTIONS"`
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
	Addresses   []string `mapstructure:"ADDRESSES"`
	MasterName  string   `mapstructure:"MASTER"`
	Password    string   `mapstructure:"PASSWORD"`
	DB          int      `mapstructure:"DB"`
	TLSDisabled bool     `mapstructure:"TLS_DISABLED"`
}

type NATS struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
}

// Load config from file
func Load(config any) error {
	v := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "_")))
	v.AutomaticEnv()
	return v.Unmarshal(config)
}
