package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type NewPostgresOptions struct {
	User         string
	Password     string
	Host         string
	Port         int
	Database     string
	Sslmode      string
	Migrate      bool
	MigrationDir string
}

func NewPostgresDb(ctx context.Context, options *NewPostgresOptions) (Database, error) {
	log.Debug().Msg("Connecting to database")
	db, err := sql.Open("postgres", formConnectionString(options.User, options.Password, options.Host, options.Database, options.Sslmode, options.Port))
	if err != nil {
		return nil, err
	}

	if options.Migrate {
		log.Debug().Msg("Running migrations")
		doMigration(db, options.MigrationDir)
	}

	log.Debug().Msg("Initializing database")
	doInit(db)

	log.Debug().Msg("Database initialized, returning connection")
	return db, nil
}

func formConnectionString(user, password, host, database, sslmode string, port int) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, password, host, port, database, sslmode)
}

func doMigration(db *sql.DB, migrationDir string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationDir,
		"postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
