package database

import (
	"errors"
	"fmt"
	"net/url"

	"campus-core/internal/config"
	"campus-core/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs database migrations
func RunMigrations(cfg *config.DatabaseConfig) error {
	migrationPath := "file://internal/database/migrations"

	// Construct migrations URL manually as golang-migrate requires URL format (postgres://)
	// whereas GORM DSN is key=value
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, url.QueryEscape(cfg.Password), cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	// Create migration instance
	m, err := migrate.New(migrationPath, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations applied successfully")
	return nil
}
