package log

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	InfoColor    string = "\033[1;32m%s\033[0m"
	WarningColor string = "\033[1;33m%s\033[0m"
	ErrorColor   string = "\033[1;31m%s\033[0m"
)

func InitLogger(appName string) {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: false}
	output.FormatLevel = func(i interface{}) string {
		if i == "warn" {
			msg := fmt.Sprintf(WarningColor, strings.ToUpper(i.(string)))
			return fmt.Sprintf("[%v]", msg)
		}

		if i == "error" {
			msg := fmt.Sprintf(ErrorColor, strings.ToUpper(i.(string)))
			return fmt.Sprintf("[%v]", msg)
		}

		msg := fmt.Sprintf(InfoColor, strings.ToUpper(i.(string)))
		return fmt.Sprintf("[%v]", msg)
	}

	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| %s |", i)
	}

	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}

	debugLog, err := strconv.ParseBool(os.Getenv("DEBUG_LOG"))
	if err != nil {
		debugLog = false
	}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger().With().Str("app", appName).Logger()
	debug := flag.Bool("debug", debugLog, "sets log level to debug")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if e := log.Debug(); e.Enabled() {
		e.Msg("Debug mode enabled")
	}
}
