package mylogger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger for easy extension
type Logger struct {
	zerolog.Logger
	Service string
}

// New creates a new Zerolog logger
func New(service, level string) (*Logger, error) {
	// Convert log level to Zerolog level
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return nil, err
	}

	// Configure Zerolog global settings
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(parsedLevel)

	// Output to console with pretty formatting
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Create the logger with the "service" context
	logger := zerolog.New(consoleWriter).
		Level(parsedLevel).
		With().
		Timestamp().
		Str("service", service).
		Logger()

	return &Logger{logger, service}, nil
}

// Action is a helper method for structured event logging
func (l *Logger) Action(event string) {
	l.Info().
		Str("action", event).
		Msg("action executed")
}
