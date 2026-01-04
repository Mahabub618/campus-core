package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port    string
	GinMode string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

// LoadConfig reads configuration from .env file and environment variables
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path + "/.env")
	viper.SetConfigType("env")

	// Read from environment variables as well
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("JWT_ACCESS_EXPIRY", "15m")
	viper.SetDefault("JWT_REFRESH_EXPIRY", "168h")
	viper.SetDefault("RATE_LIMIT_REQUESTS", 100)
	viper.SetDefault("RATE_LIMIT_DURATION", "1m")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, use environment variables and defaults
	}

	// Parse durations
	accessExpiry, err := time.ParseDuration(viper.GetString("JWT_ACCESS_EXPIRY"))
	if err != nil {
		accessExpiry = 15 * time.Minute
	}

	refreshExpiry, err := time.ParseDuration(viper.GetString("JWT_REFRESH_EXPIRY"))
	if err != nil {
		refreshExpiry = 7 * 24 * time.Hour
	}

	rateLimitDuration, err := time.ParseDuration(viper.GetString("RATE_LIMIT_DURATION"))
	if err != nil {
		rateLimitDuration = 1 * time.Minute
	}

	config := &Config{
		Server: ServerConfig{
			Port:    viper.GetString("SERVER_PORT"),
			GinMode: viper.GetString("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:        viper.GetString("JWT_SECRET"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		RateLimit: RateLimitConfig{
			Requests: viper.GetInt("RATE_LIMIT_REQUESTS"),
			Duration: rateLimitDuration,
		},
	}

	return config, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// GetRedisAddr returns the Redis address string
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
