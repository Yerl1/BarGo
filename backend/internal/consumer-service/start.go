package authservice

import (
	"context"
	"errors"
	"net/http"

	"backend/internal/auth-service/adapters/driver/myhttp"
	"backend/internal/configs"
	"backend/internal/mylogger"
)

func Execute(ctx context.Context, mylog mylogger.Logger, cfg *configs.Config) error {
	log := mylog.With().Str("Execute", "auth-service").Logger()

	server := myhttp.NewServer(ctx, mylog, cfg)

	// Run server in goroutine
	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- server.Run()
	}()

	// Wait for signal or server crash
	select {
	case <-ctx.Done():
		log.Info().Msg("Shutdown signal received")
		return server.Stop(context.Background())
	case err := <-runErrCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Msg("Server failed unexpectedly")
			return err
		}
		log.Info().Msg("Server exited normally")
		return nil
	}
}
