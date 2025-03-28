package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// LogLevel defines log severity levels
type LogLevel int

const (
	// LogLevelDebug represents debug messages
	LogLevelDebug LogLevel = iota
	// LogLevelInfo represents informational messages
	LogLevelInfo
	// LogLevelWarn represents warning messages
	LogLevelWarn
	// LogLevelError represents error messages
	LogLevelError
	// LogLevelFatal represents fatal error messages
	LogLevelFatal
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color returns ANSI color code for log level
func (l LogLevel) Color() string {
	switch l {
	case LogLevelDebug:
		return "\033[36m" // Cyan
	case LogLevelInfo:
		return "\033[32m" // Green
	case LogLevelWarn:
		return "\033[33m" // Yellow
	case LogLevelError:
		return "\033[31m" // Red
	case LogLevelFatal:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// Logger is a simple structured logger
type Logger struct {
	level      LogLevel
	out        io.Writer
	timeFormat string
	prefix     string
	useColor   bool
	mu         sync.Mutex
}

// LoggerOption is a function that configures a logger
type LoggerOption func(*Logger)

// WithLevel sets the minimum log level
func WithLevel(level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.level = level
	}
}

// WithOutput sets the output writer
func WithOutput(w io.Writer) LoggerOption {
	return func(l *Logger) {
		l.out = w
	}
}

// WithTimeFormat sets the time format
func WithTimeFormat(format string) LoggerOption {
	return func(l *Logger) {
		l.timeFormat = format
	}
}

// WithPrefix sets a prefix for all log messages
func WithPrefix(prefix string) LoggerOption {
	return func(l *Logger) {
		l.prefix = prefix
	}
}

// WithColor enables or disables colored output
func WithColor(useColor bool) LoggerOption {
	return func(l *Logger) {
		l.useColor = useColor
	}
}

// DefaultLogger is the package-level logger instance
var DefaultLogger = NewLogger()

// NewLogger creates a new logger with the given options
func NewLogger(opts ...LoggerOption) *Logger {
	// Default settings
	logger := &Logger{
		level:      LogLevelInfo,
		out:        os.Stdout,
		timeFormat: time.RFC3339,
		prefix:     "sparkci",
		useColor:   true,
	}

	// Apply options
	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

// log logs a message with the given level
func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Format args if provided
	message := msg
	if len(args) > 0 {
		message = fmt.Sprintf(msg, args...)
	}

	// Split multi-line messages
	lines := strings.Split(message, "\n")
	firstLine := lines[0]

	// Print only the message without any prefix for all log levels
	fmt.Fprintf(l.out, "%s\n", firstLine)

	// For multi-line messages, print subsequent lines without indentation
	if len(lines) > 1 {
		for _, line := range lines[1:] {
			fmt.Fprintf(l.out, "%s\n", line)
		}
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(LogLevelDebug, msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(LogLevelInfo, msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(LogLevelWarn, msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(LogLevelError, msg, args...)
}

// Variable to allow mocking os.Exit in tests
var osExit = os.Exit

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(LogLevelFatal, msg, args...)
	osExit(1)
}

// WithFields returns a new logger with the given fields attached
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	if len(fields) == 0 {
		return l
	}

	// Create field string
	var fieldParts []string
	for k, v := range fields {
		fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
	}
	fieldStr := strings.Join(fieldParts, " ")

	// Create a new logger with the prefix containing the fields
	newLogger := *l
	newLogger.prefix = fmt.Sprintf("%s [%s]", l.prefix, fieldStr)

	return &newLogger
}

// Default logger methods

// Debug logs a debug message using the default logger
func Debug(msg string, args ...interface{}) {
	DefaultLogger.Debug(msg, args...)
}

// Info logs an info message using the default logger
func Info(msg string, args ...interface{}) {
	DefaultLogger.Info(msg, args...)
}

// Warn logs a warning message using the default logger
func Warn(msg string, args ...interface{}) {
	DefaultLogger.Warn(msg, args...)
}

// Error logs an error message using the default logger
func Error(msg string, args ...interface{}) {
	DefaultLogger.Error(msg, args...)
}

// Fatal logs a fatal message and exits using the default logger
func Fatal(msg string, args ...interface{}) {
	DefaultLogger.Fatal(msg, args...)
}

// SetLogLevel sets the minimum log level for the default logger
func SetLogLevel(level LogLevel) {
	DefaultLogger.level = level
}

// GetLogLevel returns the current log level of the default logger
func GetLogLevel() LogLevel {
	return DefaultLogger.level
}

// InitLogging initializes logging based on environment variables
func InitLogging() {
	// Check for log level in environment
	logLevelEnv := strings.ToUpper(os.Getenv("SPARKCI_LOG_LEVEL"))
	if logLevelEnv != "" {
		switch logLevelEnv {
		case "DEBUG":
			SetLogLevel(LogLevelDebug)
		case "INFO":
			SetLogLevel(LogLevelInfo)
		case "WARN":
			SetLogLevel(LogLevelWarn)
		case "ERROR":
			SetLogLevel(LogLevelError)
		case "FATAL":
			SetLogLevel(LogLevelFatal)
		}
	}

	// Check if color should be disabled
	if os.Getenv("SPARKCI_NO_COLOR") == "1" || os.Getenv("NO_COLOR") == "1" {
		DefaultLogger.useColor = false
	}
}
