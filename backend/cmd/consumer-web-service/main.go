package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"backend/internal/configs"
	"backend/internal/mylogger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// -----------------------------
	// Load Configuration from YAML
	// -----------------------------
	cfg, err := configs.New("/app/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load config.yaml: %v\n", err)
		os.Exit(1)
	}

	// -----------------------------
	// Initialize Zerolog-based logger
	// -----------------------------
	logger, err := mylogger.New("consumer-web-service", "debug") // could be cfg.App.LogLevel if you add one
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to init logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info().Msg("✅ Service starting up...")
	logger.Debug().Interface("config", cfg).Msg("Loaded configuration")

	fs := http.FileServer(http.Dir("/app/web/consumer"))
	http.Handle("/", fs)

	logger.Info().Msgf("Starting consumer-web-service service on port %s", cfg.Srv.ConsumerWebServicePort)
	if err := http.ListenAndServe(":"+cfg.Srv.ConsumerWebServicePort, nil); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start consumer-web-service")
	}

	logger.Info().Msg("🛑 Shutdown signal received, cleaning up...")
	<-ctx.Done()
	logger.Info().Msg("✅ Service stopped gracefully")
}
