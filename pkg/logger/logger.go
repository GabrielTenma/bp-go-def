package logger

import (
	"fmt"
	"io"
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
func New(debug bool, broadcaster io.Writer) *Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	// Console Output (Fancy)
	consoleOutput := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	consoleOutput.FormatLevel = func(i interface{}) string {
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
	consoleOutput.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("\x1b[1m%s\x1b[0m", i)
	}

	var multi zerolog.LevelWriter
	if broadcaster != nil {
		// Broadcast Output (JSON)
		// We can't easily mix console writer and json writer in same stream if we want different formats
		// But zerolog.MultiLevelWriter allows writing to multiple writers.
		// However, zerolog writes the SAME bytes to all.
		// Trick: Use the console writer for console, and a separate logger for broadcast?
		// Actually, simpler: Let's make the main logger write JSON, and have the ConsoleWriter wrap the JSON? No.
		// Let's stick to ConsoleWriter for stdout.
		// For broadcast, we ideally want JSON.
		// To keep it simple: We will pipe the ConsoleWriter output to broadcast which is NOT JSON, but colored text.
		// Wait, user wants "live logs", usually implies reading them easily.
		// If we send colored text to web, we need to parse ANSI.
		// Better approach: Configure zerolog to write JSON to the broadcaster, and ConsoleWriter to stdout.
		// But MultiLevelWriter writes the SAME bytes.
		// Solution: Create a MultiWriter that satisfies LevelWriter, but that's complex.

		// Alternative: Just pipe the colored output to the broadcaster too. The UI can display it in a pre tag.
		// It preserves the "fancy" look in the web console too (if we strip or render ansi).
		// Let's do that for now to avoid duplicate loggers.
		multi = zerolog.MultiLevelWriter(consoleOutput, broadcaster)
	} else {
		multi = zerolog.MultiLevelWriter(consoleOutput)
	}

	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}

	z := zerolog.New(multi).Level(logLevel).With().Timestamp().Logger()

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
