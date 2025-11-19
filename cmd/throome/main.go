package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akshitmadan/throome/internal/config"
	"github.com/akshitmadan/throome/internal/logger"
	"github.com/akshitmadan/throome/pkg/gateway"
	"go.uber.org/zap"
)

var (
	// Version information (set during build)
	Version   = "0.1.0"
	BuildTime = "unknown"

	// Command-line flags
	configFile  = flag.String("config", "", "Path to configuration file")
	port        = flag.Int("port", 9000, "Server port")
	clustersDir = flag.String("clusters-dir", "./clusters", "Path to clusters directory")
	logLevel    = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	showVersion = flag.Bool("version", false, "Show version information")
)

func main() {
	flag.Parse()

	// Show version if requested
	if *showVersion {
		fmt.Printf("Throome Gateway v%s (built: %s)\n", Version, BuildTime)
		os.Exit(0)
	}

	// Initialize logger
	development := *logLevel == "debug"
	if err := logger.InitLogger(development); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Throome Gateway",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
	)

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Override with command-line flags
	if *port != 9000 {
		cfg.Server.Port = *port
	}
	if *clustersDir != "./clusters" {
		cfg.Gateway.ClustersDir = *clustersDir
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatal("Invalid configuration", zap.Error(err))
	}

	// Create gateway
	gw, err := gateway.NewGateway(cfg.Gateway.ClustersDir)
	if err != nil {
		logger.Fatal("Failed to create gateway", zap.Error(err))
	}

	// Initialize gateway
	ctx := context.Background()
	if err := gw.Initialize(ctx); err != nil {
		logger.Fatal("Failed to initialize gateway", zap.Error(err))
	}

	// Create HTTP server
	server := gateway.NewServer(cfg, gw)

	// Setup graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			logger.Error("Server error", zap.Error(err))
		}
	}()

	logger.Info("Throome Gateway is running",
		zap.Int("port", cfg.Server.Port),
		zap.String("clusters_dir", cfg.Gateway.ClustersDir),
	)

	// Wait for shutdown signal
	<-shutdown

	logger.Info("Shutdown signal received, gracefully shutting down...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}

	// Shutdown gateway
	if err := gw.Shutdown(shutdownCtx); err != nil {
		logger.Error("Gateway shutdown error", zap.Error(err))
	}

	logger.Info("Throome Gateway stopped")
}

// loadConfig loads the application configuration
func loadConfig() (*config.AppConfig, error) {
	if *configFile != "" {
		return config.LoadConfig(*configFile)
	}

	// Return default configuration
	return config.DefaultConfig(), nil
}
