package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func New(level string) *Logger {
	var l zerolog.Level

	switch level {
	case "debug":
		l = zerolog.DebugLevel
	case "info":
		l = zerolog.InfoLevel
	case "error":
		l = zerolog.ErrorLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)
	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info().Msgf(message, args...)
}
func (l *Logger) Error(err error, message string, args ...interface{}) {
	l.logger.Error().Err(err).Msgf(message, args...)

}
