package gotrunks

import (
	"io"

	llog "github.com/sirupsen/logrus"
)

// Level is an alias type to the logger implementaion
type Level llog.Level

// Hook :nodoc:
type Hook llog.Hook

// Format :nodoc:
type Format uint8

// std logger instance
var (
	std = llog.New()
	env = "unknown"
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = Level(llog.PanicLevel)
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = Level(llog.FatalLevel)
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = Level(llog.ErrorLevel)
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = Level(llog.WarnLevel)
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = Level(llog.InfoLevel)
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = Level(llog.DebugLevel)
)

const (
	_ Format = iota
	// JSONFormat iota
	JSONFormat
	// TextFormat iota
	TextFormat
)

// Logger is an interface for general logging
type Logger interface {
	Print(...interface{})
	Println(...interface{})
	Printf(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

// Standard :nodoc:
func Standard() Logger {
	return std
}

// GetLogger is a function to get default logger
func GetLogger() Logger {
	return std.WithField("env", env)
}

// SetOutput to change logger argsput
func SetOutput(w io.Writer) {
	std.Out = w
}

// SetLevel of the logger
func SetLevel(level Level) {
	std.Level = llog.Level(level)
}

// SetEnv :nodoc:
func SetEnv(e string) {
	env = e
}

// SetFormat for the logger
func SetFormat(format Format) {
	switch format {
	case JSONFormat:
		std.Formatter = &llog.JSONFormatter{}
	default:
		std.Formatter = &llog.TextFormatter{}
	}
}

// AddHook to Standard Logger
func AddHook(h Hook) {
	std.Hooks.Add(h)
}

// From adds package and function name where log funcs are called
func From(pkg, fn string) Logger {
	return std.WithFields(llog.Fields{
		"package":  pkg,
		"function": fn,
	})
}

// Print is an alias method to the logger implementaion
func Print(args ...interface{}) {
	GetLogger().Print(args...)
}

// Printf is an alias method to the logger implementaion
func Printf(format string, args ...interface{}) {
	GetLogger().Printf(format, args...)
}

// Debug is an alias method to the logger implementaion
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf is an alias method to the logger implementaion
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info is an alias method to the logger implementaion
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof is an alias method to the logger implementaion
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn is an alias method to the logger implementaion
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf is an alias method to the logger implementaion
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error is an alias method to the logger implementaion
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf is an alias method to the logger implementaion
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal is an alias method to the logger implementaion
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf is an alias method to the logger implementaion
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}
