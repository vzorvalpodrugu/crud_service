package config

import (
	"fmt"
	"os"
)

type Config struct {
	app AppConfig
	db  DBConfig
}

type AppConfig struct {
	port string
	host string
}

type DBConfig struct {
	dbname   string
	host     string
	port     string
	user     string
	password string
}

// DSN
func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"dbname=%s host=%s post=%s user=%s password=%s",
	)
}

func Load() (*Config, error) {
	cfg := &Config{
		app: AppConfig{
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_HOST"),
		},
		db: DBConfig{
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
		},
	}

	return cfg, nil
}
