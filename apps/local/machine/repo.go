package machine

import (
	"context"
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

type Repo interface {
	GetMachineById(ctx context.Context, id string) (*Model, error)
	GetMachines(ctx context.Context, filters *GetMachinesFilters) ([]*Model, error)
	CreateMachine(ctx context.Context, m *Model) error
	UpdateMachine(ctx context.Context, m *Model) error
	// DeleteMachine(ctx context.Context, id string) error
	AcknowledgeHeartbeat(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status MachineStatus) error
}

func CreatePostgresConnectionString(config RepoConfig) string {
	sslMode := "disable"
	if config.UseSsl {
		sslMode = "enable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, sslMode)
}
