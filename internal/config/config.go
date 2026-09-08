package config

import (
	"fmt"
	"os"
)

type Config struct {
	App   AppConfig
	Db    DBConfig
	Redis RedisConfig
}
type RedisConfig struct {
	Host string
	Port string
}

type AppConfig struct {
	Port string
	host string
}

type DBConfig struct {
	Dbname   string
	Host     string
	Port     string
	User     string
	Password string
}

// DSN
func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"dbname=%s host=%s port=%s user=%s password=%s sslmode=disable",
		db.Dbname, db.Host, db.Port, db.User, db.Password,
	)
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			os.Getenv("APP_PORT"),
			os.Getenv("POSTGRES_HOST"),
		},
		Db: DBConfig{
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
		},
		Redis: RedisConfig{
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		},
	}

	return cfg, nil
}
