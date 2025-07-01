package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger = logger.Hook(&FileNameHook{})

	log.Logger = logger

	return nil
}

type FileNameHook struct{}

func (_ *FileNameHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	e = e.Caller(3)
}

type ConfigureLoggerOptions struct {
	Level zerolog.Level
}

type LogConfig func(*ConfigureLoggerOptions)

func WithLevel(l zerolog.Level) LogConfig {
	return func(c *ConfigureLoggerOptions) {
		c.Level = l
	}
}
