// Package logging provides logging tools
package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigureLogger initializes zerolog's logger.
func ConfigureLogger(opts ...LogConfig) error {
	config := ConfigureLoggerOptions{
		Level: zerolog.InfoLevel,
	}

	for _, opt := range opts {
		opt(&config)
	}

	zerolog.SetGlobalLevel(config.Level)

	zerolog.FormattedLevels[zerolog.TraceLevel] = "TRACE  "
	zerolog.FormattedLevels[zerolog.DebugLevel] = "DEBUG  "
	zerolog.FormattedLevels[zerolog.InfoLevel] = "INFO   "
	zerolog.FormattedLevels[zerolog.WarnLevel] = "WARNING"
	zerolog.FormattedLevels[zerolog.ErrorLevel] = "ERROR  "
	zerolog.FormattedLevels[zerolog.FatalLevel] = "FATAL  "
	zerolog.FormattedLevels[zerolog.PanicLevel] = "PANIC  "

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "timestamp"

	logger := log.Output(
		zerolog.ConsoleWriter{ //nolint:exhaustruct  // Other fields are Zero by design
			Out: os.Stderr,
		},
	)
	logger = logger.Hook(&FileNameHook{})

	log.Logger = logger

	return nil
}

const fileNameHookFrameSkip = 3

// FileNameHook is a hook to add the name of the file from which a log was emitted.
type FileNameHook struct{}

// Run applies the FileNameHook.
func (*FileNameHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Caller(fileNameHookFrameSkip)
}

// ConfigureLoggerOptions is the configuration for our logger.
type ConfigureLoggerOptions struct {
	Level zerolog.Level
}

// LogConfig is a function to change the logger configuration.
type LogConfig func(*ConfigureLoggerOptions)

// WithLevel changes the log level.
func WithLevel(l zerolog.Level) LogConfig {
	return func(c *ConfigureLoggerOptions) {
		c.Level = l
	}
}
