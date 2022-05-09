package utils

import (
	"errors"
	"fmt"
	"io"
)

type LogLevel int

const (
	Info LogLevel = iota
	Warn
	Error
	Ignore
)

// Logger act as a regular logger that counts logged errors and warnings.
type Logger interface {
	Log(LogLevel, any)
	Logf(LogLevel, string, ...any)
	LogIfNotNil(LogLevel, any) bool
	LogIfNotNilf(LogLevel, any, string, ...any) bool
	Info(any)
	Infof(string, ...any)
	Warn(any)
	Warnf(string, ...any)
	Error(any)
	Errorf(string, ...any)
	ErrorIfNotNil(any) bool
	ErrorIfNotNilf(any, string, ...any) bool
	WarnIfNotNil(any) bool
	WarnIfNotNilf(any, string, ...any) bool
	Counters() LogCounters
	Warnings() int
	Errors() int
	MustHaveNoWarnings() error
	MustHaveNoErrors() error
	MustHaveNoWarningsf(string, ...any) error
	MustHaveNoErrorsf(string, ...any) error
	AppendCounters(LogCounters)
	// NewCounter creates a cloned logger with the same output and log level but with new counters
	NewCounter() Logger
	WithLevel(level LogLevel)
}

type LogCounters struct {
	Errors   int
	Warnings int
}

func NewLogger(out io.Writer) Logger {
	return &logger{
		output: out,
	}
}
func NewLoggerWithLevel(out io.Writer, level LogLevel) Logger {
	return &logger{
		output: out,
		level:  level,
	}
}

type logger struct {
	output io.Writer
	level  LogLevel
	LogCounters
}

func (l *logger) WithLevel(level LogLevel)          { l.level = level }
func (l *logger) Info(arg any)                      { l.Log(Info, arg) }
func (l *logger) Infof(format string, args ...any)  { l.Logf(Info, format, args...) }
func (l *logger) Warn(arg any)                      { l.Log(Warn, arg) }
func (l *logger) Warnf(format string, args ...any)  { l.Logf(Warn, format, args...) }
func (l *logger) Error(arg any)                     { l.Log(Error, arg) }
func (l *logger) Errorf(format string, args ...any) { l.Logf(Error, format, args...) }
func (l *logger) Counters() LogCounters             { return l.LogCounters }
func (l *logger) Warnings() int                     { return l.LogCounters.Warnings }
func (l *logger) Errors() int                       { return l.LogCounters.Errors }
func (l *logger) ErrorIfNotNil(arg any) bool        { return l.LogIfNotNil(Error, arg) }
func (l *logger) WarnIfNotNil(arg any) bool         { return l.LogIfNotNil(Warn, arg) }
func (l *logger) WarnIfNotNilf(arg any, format string, args ...any) bool {
	return l.LogIfNotNilf(Warn, arg, format, args...)
}
func (l *logger) ErrorIfNotNilf(arg any, format string, args ...any) bool {
	return l.LogIfNotNilf(Error, arg, format, args...)
}

func (l *logger) AppendCounters(counters LogCounters) {
	l.LogCounters.Errors += counters.Errors
	l.LogCounters.Warnings += counters.Warnings
}

func (l *logger) LogIfNotNil(action LogLevel, err any) bool {
	if err == nil {
		return false
	}
	l.Log(action, err)
	return true
}
func (l *logger) LogIfNotNilf(action LogLevel, arg any, format string, args ...any) bool {
	if arg == nil {
		return false
	}
	l.Log(action, fmt.Sprintf(format, args...))
	return true
}

func (l *logger) Logf(level LogLevel, format string, args ...any) {
	l.Log(level, fmt.Sprintf(format, args...))
}
func (l *logger) Log(level LogLevel, arg any) {
	write := func(string) {}
	if l.level <= level {
		write = func(levelStr string) { fmt.Fprintln(l.output, levelStr, arg) }
	}
	switch level {
	case Info:
		write("[Info]")
	case Warn:
		write("[Warning]")
		l.LogCounters.Warnings++
	case Error:
		write("[Error]")
		l.LogCounters.Errors++
	}
}

func (l *logger) MustHaveNoErrors() error {
	if l.LogCounters.Errors == 0 {
		return nil
	}
	return errors.New(mustHaveCountMessage(l.LogCounters.Errors, "error", "errors", "") + " logged")
}

func (l *logger) NewCounter() Logger {
	return &logger{output: l.output, level: l.level}
}

func (l *logger) MustHaveNoWarnings() error {
	if l.LogCounters.Errors == 0 && l.LogCounters.Warnings == 0 {
		return nil
	}
	errMessage := mustHaveCountMessage(l.LogCounters.Errors, "error", "errors", "")
	warnMessage := mustHaveCountMessage(l.LogCounters.Warnings, "warning", "warnings", " and ")
	return errors.New(errMessage + warnMessage + " logged")
}

func (l *logger) MustHaveNoWarningsf(format string, args ...any) error {
	if l.MustHaveNoWarnings() != nil {
		return fmt.Errorf(format, args...)
	}
	return nil
}
func (l *logger) MustHaveNoErrorsf(format string, args ...any) error {
	if l.MustHaveNoErrors() != nil {
		return fmt.Errorf(format, args...)
	}
	return nil
}

func mustHaveCountMessage(count int, singular string, plural string, prefix string) string {
	switch {
	case count == 1:
		return fmt.Sprintf("%sone %s", prefix, singular)
	case count > 1:
		return fmt.Sprintf("%s%d %s", prefix, count, plural)
	}
	return ""
}