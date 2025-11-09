package myhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"backend/internal/business-service/adapters/driven/db"
	"backend/internal/business-service/adapters/driver/myhttp/handlers"
	"backend/internal/business-service/adapters/driver/myhttp/middleware"
	"backend/internal/business-service/core/services"
	"backend/internal/configs"
	"backend/internal/mylogger"
)

var ErrServerClosed = errors.New("Server closed")

const WaitTime = 10

type Server struct {
	mux    *http.ServeMux
	cfg    *configs.Config
	srv    *http.Server
	mylog  mylogger.Logger
	db     *db.DB
	ctx    context.Context
	appCtx context.Context
	mu     sync.Mutex
	wg     sync.WaitGroup
	mdl    *middleware.Middleware
}

func NewServer(ctx context.Context, mylog mylogger.Logger, cfg *configs.Config) *Server {
	return &Server{
		ctx:    ctx,
		appCtx: context.Background(),
		cfg:    cfg,
		mylog:  mylog,
		mux:    http.NewServeMux(),
	}
}

// Run initializes routes and starts listening. It returns when the server stops.
func (s *Server) Run() error {
	log := s.mylog.With().Str("action", "Run").Logger()

	log.Info().Msg("server_started")

	// Initialize database connection
	if err := s.initializeDatabase(); err != nil {
		log.Err(err).Msg("Failed to connect to database")
		return err
	}
	log.Info().Msg("Successful database connection")

	mdl := middleware.New(s.ctx, s.cfg.App.JwtSecret, s.mylog)
	s.mdl = mdl
	// Configure routes and handlers
	s.Configure()

	log.Info().Str("port", s.cfg.Srv.BusinessServicePort).Msg("this port")

	s.mu.Lock()
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%v", s.cfg.Srv.BusinessServicePort),
		Handler: s.mdl.CorsMiddleware(s.mux),
	}
	s.mu.Unlock()

	log.Info().Msg("server is running")
	// Start the HTTP server and handle graceful shutdown
	return s.startHTTPServer()
}

// Stop provides a programmatic shutdown. Accepts a context for timeout control.
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log := s.mylog.With().Str("action", "Run").Logger()

	log.Info().Msg("Shutting down HTTP server...")

	s.wg.Wait()

	if s.srv != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, WaitTime*time.Second)
		defer cancel()

		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			log.Err(err).Msg("Failed to shut down HTTP server gracefully")
			return fmt.Errorf("http server shutdown: %w", err)
		}
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.Err(err).Msg("Failed to close database")
			return fmt.Errorf("db close: %w", err)
		}
		log.Info().Msg("Database closed")
	}

	log.Info().Msg("HTTP server shut down gracefully")

	return nil
}

func (s *Server) startHTTPServer() error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-s.ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

// Configure sets up the HTTP handlers for various APIs including Market Data, Data Mode control, and Health checks.
func (s *Server) Configure() {
	businessRepo := db.NewBusinessRepo(s.ctx, s.db, s.mylog)

	businessService := services.NewBusinessService(s.ctx, businessRepo, s.cfg.App.JwtSecret, s.mylog)

	businessHandler := handlers.NewBusinessHandler(s.ctx, businessService, s.mylog)

	s.mux.Handle("GET /consumer/health", businessHandler.HealthHandler())

	s.mux.Handle("GET /stores", s.mdl.Wrap(businessHandler.GetStores()))
	s.mux.Handle("POST /stores", s.mdl.Wrap(businessHandler.AddStore()))
	s.mux.Handle("PUT /stores", s.mdl.Wrap(businessHandler.UpdateStore()))
	s.mux.Handle("DELETE /stores", s.mdl.Wrap(businessHandler.DeleteStore()))

	s.mux.Handle("GET /products", s.mdl.Wrap(businessHandler.GetProducts()))
	s.mux.Handle("POST /products", s.mdl.Wrap(businessHandler.AddProduct()))
	s.mux.Handle("PUT /products", s.mdl.Wrap(businessHandler.UpdateProduct()))
	s.mux.Handle("DELETE /products", s.mdl.Wrap(businessHandler.DeleteProduct()))
}

func (s *Server) initializeDatabase() error {
	db, err := db.Start(s.ctx, s.cfg.DB, s.mylog)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	s.db = db
	return nil
}
