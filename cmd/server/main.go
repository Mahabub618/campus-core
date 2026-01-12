package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"campus-core/internal/config"
	"campus-core/internal/database"
	"campus-core/internal/router"
	"campus-core/internal/utils"
	"campus-core/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Server.GinMode); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Campus Core Server",
		zap.String("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.GinMode),
	)

	// Initialize validator
	if err := utils.InitValidator(); err != nil {
		logger.Fatal("Failed to initialize validator", zap.Error(err))
	}

	// Connect to database
	db, err := database.ConnectDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.CloseDB()

	// Run database migrations
	if err := database.RunMigrations(&cfg.Database); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Seed database
	seeder := database.NewSeeder(db)
	if err := seeder.SeedAll(); err != nil {
		logger.Error("Failed to seed database", zap.Error(err))
	}

	// Connect to Redis (optional, continue if fails)
	_, err = database.ConnectRedis(&cfg.Redis)
	if err != nil {
		logger.Warn("Failed to connect to Redis, rate limiting will be disabled", zap.Error(err))
	} else {
		defer database.CloseRedis()
	}

	// Setup router
	r := router.NewRouter(cfg, db)
	engine := r.Setup()

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		logger.Info("Server listening", zap.String("address", addr))
		if err := engine.Run(addr); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	logger.Info("Server exited gracefully")
}
