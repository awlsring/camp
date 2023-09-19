package repo

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/apps/local/machine"
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
	GetMachineById(ctx context.Context, id string) (*machine.Model, error)
	GetMachines(ctx context.Context, filters *GetMachinesFilters) ([]*machine.Model, error)
	CreateMachine(ctx context.Context, m *machine.Model) error
	UpdateMachine(ctx context.Context, m *machine.Model) error
	// DeleteMachine(ctx context.Context, id string) error
	AcknowledgeHeartbeat(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status machine.MachineStatus) error
}

func CreatePostgresConnectionString(config RepoConfig) string {
	sslMode := "disable"
	if config.UseSsl {
		sslMode = "enable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, sslMode)
}
