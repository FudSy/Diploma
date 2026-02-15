package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level string `default:"error"`
}

func (c Config) Init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC1123,
	})

	level := zerolog.InfoLevel
	
	switch c.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn", "warning":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	case "panic":
		level = zerolog.PanicLevel
	default:
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	log.Info().Str("level", level.String()).Msg("Логгер инициализирован")
}