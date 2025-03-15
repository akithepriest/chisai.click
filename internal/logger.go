package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type ChisaiLogger struct {
	zerolog.Logger
}

var Logger ChisaiLogger

// Custom logger implementation
func NewLogger() ChisaiLogger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}

	zerolog := zerolog.New(output).With().Caller().Timestamp().Logger()
	Logger = ChisaiLogger{zerolog}
	return Logger
}

func (l *ChisaiLogger) LogInfo() *zerolog.Event {
	return l.Logger.Info()
}

func (l *ChisaiLogger) LogError() *zerolog.Event {
	return l.Logger.Error()
}

func (l *ChisaiLogger) LogDebug() *zerolog.Event {
	return l.Logger.Debug()
}

func (l *ChisaiLogger) LogWarn() *zerolog.Event {
	return l.Logger.Warn()
}

func (l *ChisaiLogger) LogFatal() *zerolog.Event {
	return l.Logger.Fatal()
}
