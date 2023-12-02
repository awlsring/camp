package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/awlsring/camp/internal/pkg/database"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v2"
)

type MetricsConfig struct {
	ServiceName string `json:"service_name" yaml:"service_name" toml:"service_name"`
	Address     string `json:"address" yaml:"address" toml:"address"`
}

type ServerConfig struct {
	Address   string   `json:"address" yaml:"address" toml:"address"`
	ApiKeys   []string `json:"api_keys" yaml:"api_keys" toml:"api_keys"`
	AgentKeys []string `json:"agent_keys" yaml:"agent_keys" toml:"agent_keys"`
}

type DatabaseConfig struct {
	Driver   string                   `json:"driver" yaml:"driver" toml:"driver"`
	Postgres *database.PostgresConfig `json:"postgres" yaml:"postgres" toml:"postgres"`
}

type Config struct {
	Server   ServerConfig   `json:"server" yaml:"server" toml:"server"`
	Metrics  MetricsConfig  `json:"metrics" yaml:"metrics" toml:"metrics"`
	Database DatabaseConfig `json:"database" yaml:"database" toml:"database"`
}

func LoadConfig(ctx context.Context, path string) (*Config, error) {
	log := logger.FromContext(ctx)
	log.Info().Msgf("Loading config from %s", path)

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file")
	}

	var cfg *Config
	switch true {
	case strings.HasSuffix(path, ".json"):
		log.Debug().Msg("Loading JSON config")
		cfg, err = loadJsonConfig(file)
	case strings.HasSuffix(path, ".toml"):
		log.Debug().Msg("Loading TOML config")
		cfg, err = loadTomlConfig(file)
	case strings.HasSuffix(path, ".yaml"):
		log.Debug().Msg("Loading YAML config")
		cfg, err = loadYamlConfig(file)
	default:
		return nil, fmt.Errorf("unknown config file type, %s", path)
	}

	// validate / set defaults
	log.Debug().Msg("Validating config")
	setConfigDefaults(cfg)

	return cfg, err
}

func loadYamlConfig(file []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadTomlConfig(file []byte) (*Config, error) {
	var config Config
	if err := toml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadJsonConfig(file []byte) (*Config, error) {
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setConfigDefaults(cfg *Config) {
	setDefaultMetricsConfig(cfg)
	setDefaultServerConfig(cfg)
	setDefaultDatabaseConfig(cfg)
}

func setDefaultMetricsConfig(cfg *Config) {
	if cfg.Metrics.Address == "" {
		cfg.Metrics.Address = "127.0.0.1:9032"
	}
	if cfg.Metrics.ServiceName == "" {
		cfg.Metrics.ServiceName = "CampLocal"
	}
}

func setDefaultServerConfig(cfg *Config) {
	if cfg.Server.Address == "" {
		cfg.Server.Address = "127.0.0.1:8032"
	}
}

func setDefaultDatabaseConfig(cfg *Config) {
	if cfg.Database.Driver == "" {
		cfg.Database.Driver = "postgres"
	}
	if cfg.Database.Driver == "postgres" {
		if cfg.Database.Postgres == nil {
			cfg.Database.Postgres = &database.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "postgres",
				Database: "camplocal",
				UseSsl:   false,
			}
		}
		if cfg.Database.Postgres.Host == "" {
			cfg.Database.Postgres.Host = "localhost"
		}
		if cfg.Database.Postgres.Port == 0 {
			cfg.Database.Postgres.Port = 5432
		}
		if cfg.Database.Postgres.Username == "" {
			cfg.Database.Postgres.Username = "postgres"
		}
		if cfg.Database.Postgres.Password == "" {
			cfg.Database.Postgres.Password = "postgres"
		}
		if cfg.Database.Postgres.Database == "" {
			cfg.Database.Postgres.Database = "camplocal"
		}
		if cfg.Database.Postgres.UseSsl == false {
			cfg.Database.Postgres.UseSsl = false
		}
	}
}
