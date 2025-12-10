package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps the zerolog logger
type Logger struct {
	z zerolog.Logger
}

// New creates a new fancy logger
func New(debug bool) *Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	output.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				l = "\x1b[36mDEBUG 🐛\x1b[0m"
			case "info":
				l = "\x1b[32mINFO 🚀\x1b[0m"
			case "warn":
				l = "\x1b[33mWARN ⚠️\x1b[0m"
			case "error":
				l = "\x1b[31mERROR ❌\x1b[0m"
			case "fatal":
				l = "\x1b[31mFATAL 💀\x1b[0m"
			case "panic":
				l = "\x1b[31mPANIC 💥\x1b[0m"
			default:
				l = strings.ToUpper(ll)
			}
		} else {
			if i == nil {
				l = strings.ToUpper(fmt.Sprintf("%s", i))
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))
			}
		}
		return l
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("\x1b[1m%s\x1b[0m", i)
	}

	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}

	z := zerolog.New(output).Level(logLevel).With().Timestamp().Logger()

	return &Logger{z: z}
}

// Info logs an info message
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.log(l.z.Info(), msg, keyvals...)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, keyvals ...interface{}) {
	if err != nil {
		l.z.Error().Err(err).Fields(keyvals).Msg(msg)
	} else {
		l.log(l.z.Error(), msg, keyvals...)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.log(l.z.Debug(), msg, keyvals...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.log(l.z.Warn(), msg, keyvals...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error) {
	if err != nil {
		l.z.Fatal().Err(err).Msg(msg)
	} else {
		l.z.Fatal().Msg(msg)
	}
}

func (l *Logger) log(e *zerolog.Event, msg string, keyvals ...interface{}) {
	if len(keyvals)%2 != 0 {
		e.Msg(msg + " (odd number of keyvals caused metadata drop)")
		return
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keyvals[i])
		}
		e.Interface(key, keyvals[i+1])
	}
	e.Msg(msg)
}
