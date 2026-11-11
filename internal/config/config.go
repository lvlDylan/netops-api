package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port    string `envconfig:"PORT" default:"8080"`
	LogFile bool   `envconfig:"LOGFILE" default:"false"`
	Verbose bool   `envconfig:"VERBOSE" default:"false"`

	Database string `envconfig:"DATABASE" default:"mariadb"`
	DBUser   string `envconfig:"DB_USER" default:"mariadb"`
	DBPass   string `envconfig:"DB_PASS" default:"mariadb"`
	DBHost   string `envconfig:"DB_HOST" default:"localhost"`
	DBPort   string `envconfig:"DB_PORT" default:"3306"`
	DBName   string `envconfig:"DB_NAME" default:"mariadb"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, fmt.Errorf("impossible de charger la configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
