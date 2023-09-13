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

func createPostgresConnectionString(config RepoConfig) string {
	sslMode := "disable"
	if config.UseSsl {
		sslMode = "enable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, sslMode)
}

// type PostgresRepo struct {
// 	database database.Database
// }

// func NewRepo(db database.Database) Repo {
// 	return &PostgresRepo{
// 		database: db,
// 	}
// }

// func (r *PostgresRepo) GetMachine(ctx context.Context, id string) (*Model, error) {
// 	var m Model
// 	err := r.database.QueryRowContext(ctx, "SELECT * FROM machines WHERE id = $1", id).Scan(&m.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &m, nil
// }

// func (r *PostgresRepo) GetMachines(ctx context.Context) ([]*Model, error) {
// 	rows, err := r.database.QueryContext(ctx, "SELECT * FROM machines")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var machines []*Model
// 	for rows.Next() {
// 		var m Model
// 		err := rows.Scan(&m.Id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		machines = append(machines, &m)
// 	}

// 	return machines, nil
// }

// func (r *PostgresRepo) CreateMachine(ctx context.Context, m *Model) error {
// 	_, err := r.database.ExecContext(ctx, "INSERT INTO machines (id) VALUES ($1)", m.Id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *PostgresRepo) UpdateMachine(ctx context.Context, m *Model) error {
// 	_, err := r.database.ExecContext(ctx, "UPDATE machines SET id = $1 WHERE id = $2", m.Id, m.Id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *PostgresRepo) DeleteMachine(ctx context.Context, id string) error {
// 	_, err := r.database.ExecContext(ctx, "DELETE FROM machines WHERE id = $1", id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
