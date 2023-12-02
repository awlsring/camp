package database

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string `json:"host" yaml:"host" toml:"host"`
	Port     int    `json:"port" yaml:"port" toml:"port"`
	Username string `json:"username" yaml:"username" toml:"username"`
	Password string `json:"password" yaml:"password" toml:"password"`
	Database string `json:"database" yaml:"database" toml:"database"`
	UseSsl   bool   `json:"use_ssl" yaml:"use_ssl" toml:"use_ssl"`
}

func CreatePostgresConnectionString(config PostgresConfig) string {
	sslMode := "disable"
	if config.UseSsl {
		sslMode = "enable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, sslMode)
}
