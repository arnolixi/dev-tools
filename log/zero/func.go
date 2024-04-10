package zero

import (
	"fmt"
	"github.com/go-stack/stack"
	"os"
)

func Debugf(format string, args ...interface{}) {
	logger.Debug().Msgf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Info().Msgf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warn().Msgf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	frames := stack.Trace()
	errMsg := fmt.Sprintf(format, args...)
	for _, frame := range frames {
		frameStr := fmt.Sprintf("%+v", frame)
		errMsg += fmt.Sprintf("\n\t\t\t%v -> %n()", frameStr, frame)
	}
	logger.Error().Msgf(errMsg)
	os.Exit(1)
}
