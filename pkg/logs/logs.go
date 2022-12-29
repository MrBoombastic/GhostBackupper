package logs

import (
	"github.com/BOOMfinity/golog"
)

var logger golog.Logger

func init() {
	logger = golog.New("GhostBackupper")
}

func Info(s string) {
	logger.Info().Send(s)
}
func Error(s string) {
	logger.Error().Send(s)
}
