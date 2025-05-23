package godb

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/mickamy/godb/config"
)

var (
	ErrMigrateNoChange = errors.New("no change")
)

func Migrate(cfg config.Config) error {
	dbURL, err := cfg.Database.URL()
	if err != nil {
		return fmt.Errorf("failed to get database URL: %w", err)
	}

	migrationPath, err := filepath.Abs(cfg.Migrations.Dir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for migrations: %w", err)
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), dbURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrMigrateNoChange
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
