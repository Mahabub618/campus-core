package database

import (
	"fmt"

	"campus-core/internal/config"
	"campus-core/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB establishes a connection to PostgreSQL database
func ConnectDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	logger.Info("Database connected successfully", zap.String("host", cfg.Host), zap.String("database", cfg.DBName))

	return db, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate runs auto-migration for all models
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}
	logger.Info("Database migration completed successfully")
	return nil
}
