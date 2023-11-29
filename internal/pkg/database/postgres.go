package database

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	UseSsl   bool
}

func CreatePostgresConnectionString(config PostgresConfig) string {
	sslMode := "disable"
	if config.UseSsl {
		sslMode = "enable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, sslMode)
}
