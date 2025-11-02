package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"backend/internal/auth-service/core/myerrors"
	"backend/internal/configs"
	"backend/internal/mylogger"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	ctx   context.Context
	cfg   *configs.DBconfig
	mylog mylogger.Logger
	conn  *pgx.Conn
	mu    *sync.Mutex
}

// Start initializes and returns a new DB instance with a single connection
func Start(ctx context.Context, dbCfg *configs.DBconfig, mylog mylogger.Logger) (*DB, error) {
	log := mylog.With().Str("db", "Start").Logger()
	d := &DB{
		cfg:   dbCfg,
		ctx:   ctx,
		mylog: mylog,
		mu:    &sync.Mutex{},
	}

	// Attempt to connect with retries
	if err := d.connect(); err != nil {
		log.Err(err).Msg("Failed to establish DB connection")
		return nil, err
	}

	log.Info().Msg("DB connection established successfully")
	return d, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	log := d.mylog.With().Str("db", "Close").Logger()

	if err := d.conn.Close(d.ctx); err != nil {
		log.Err(err).Msg("Failed to close DB connection")
		return fmt.Errorf("close database connection: %v", err)
	}
	log.Info().Msg("Database connection closed successfully")
	return nil
}

// IsAlive checks if the DB connection is still alive
func (d *DB) IsAlive() error {
	log := d.mylog.With().Str("db", "IsAlive").Logger()

	if d.conn == nil {
		return fmt.Errorf("DB is not initialized")
	}

	// Ping the DB to ensure the connection is alive
	if err := d.conn.Ping(d.ctx); err != nil {
		// If the ping fails, try reconnecting
		log.Err(err).Msg("DB ping failed, attempting to reconnect...")
		if connectionErr := d.connect(); connectionErr != nil {
			return myerrors.ErrDBConnClosed
		}
	}

	return nil
}

// connect establishes a new connection with retry logic
func (d *DB) connect() error {
	log := d.mylog.With().Str("db", "connect").Logger()

	var lastErr error
	for i := 0; i < d.cfg.MaxRetries; i++ {
		// Build connection string
		connStr := fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v?sslmode=disable",
			d.cfg.User,
			d.cfg.Password,
			d.cfg.Host,
			d.cfg.Port,
			d.cfg.Database,
		)

		// Attempt to establish a connection
		conn, err := pgx.Connect(d.ctx, connStr)
		if err != nil {
			// Save the error and retry

			lastErr = fmt.Errorf("failed to connect to database: %w", err)
			log.Err(err).Msg(fmt.Sprintf("DB connection attempt %d failed", i+1))

			// Exponential backoff (1s, 2s, 3s, etc.)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		// Successfully connected
		d.conn = conn
		log.Info().Msg("Successfully connected to the database")
		return nil
	}

	// Return the last error after all retries have failed
	return fmt.Errorf("failed to connect to the database after %d attempts: %w", d.cfg.MaxRetries, lastErr)
}
