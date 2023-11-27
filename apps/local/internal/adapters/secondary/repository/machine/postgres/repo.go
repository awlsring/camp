package postgres

import (
	"fmt"
)

type GetMachinesFilters struct {
}

type RepoConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	UseSsl   bool
}

func CreatePostgresConnectionString(config RepoConfig) string {
	sslMode := "disable"
	if config.UseSsl {
		sslMode = "enable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, sslMode)
}
