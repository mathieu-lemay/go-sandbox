package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func ConfigureLogger(opts ...LogConfig) error {
	config := ConfigureLoggerOptions{
		Level: zerolog.InfoLevel,
	}

	for _, opt := range opts {
		opt(&config)
	}

	zerolog.SetGlobalLevel(config.Level)

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
