package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevelFatal, "FATAL"},
		{LogLevel(99), "UNKNOWN"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.level.String(); got != test.expected {
				t.Errorf("LogLevel.String() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestLoggerOutput(t *testing.T) {
	var buf bytes.Buffer

	// Create logger with buffer output and no color
	logger := NewLogger(
		WithOutput(&buf),
		WithColor(false),
		WithPrefix("test"),
	)

	// Set fixed time format without using actual time to make tests deterministic
	logger.timeFormat = ""

	// Test different log levels
	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		level    LogLevel
		message  string
		expected string
	}{
		{
			name:     "Debug",
			logFunc:  logger.Debug,
			level:    LogLevelDebug,
			message:  "debug message",
			expected: "test [DEBUG] debug message",
		},
		{
			name:     "Info",
			logFunc:  logger.Info,
			level:    LogLevelInfo,
			message:  "info message",
			expected: "test [INFO] info message",
		},
		{
			name:     "Warn",
			logFunc:  logger.Warn,
			level:    LogLevelWarn,
			message:  "warn message",
			expected: "test [WARN] warn message",
		},
		{
			name:     "Error",
			logFunc:  logger.Error,
			level:    LogLevelError,
			message:  "error message",
			expected: "test [ERROR] error message",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Clear buffer
			buf.Reset()

			// Set log level to test level
			logger.level = test.level

			// Log message
			test.logFunc(test.message)

			// Check output
			got := strings.TrimSpace(buf.String())
			if !strings.Contains(got, test.expected) {
				t.Errorf("Expected log output to contain %q, got %q", test.expected, got)
			}
		})
	}
}

func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer

	// Create logger with buffer output and no color
	logger := NewLogger(
		WithOutput(&buf),
		WithColor(false),
		WithPrefix("test"),
	)

	// Set empty time format for deterministic tests
	logger.timeFormat = ""

	// Test WithFields
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	fieldLogger := logger.WithFields(fields)
	fieldLogger.Info("field test")

	got := strings.TrimSpace(buf.String())

	// Check that both fields are in the output
	if !strings.Contains(got, "key1=value1") || !strings.Contains(got, "key2=42") {
		t.Errorf("Expected log output to contain fields, got %q", got)
	}
}

func TestLoggerMultilineMessage(t *testing.T) {
	var buf bytes.Buffer

	// Create logger with buffer output and no color
	logger := NewLogger(
		WithOutput(&buf),
		WithColor(false),
		WithPrefix("test"),
	)

	// Set empty time format for deterministic tests
	logger.timeFormat = ""

	// Test multiline message
	logger.Info("line1\nline2\nline3")

	lines := strings.Split(buf.String(), "\n")
	if len(lines) < 3 {
		t.Fatalf("Expected at least 3 lines, got %d", len(lines))
	}

	// First line should contain "line1"
	if !strings.Contains(lines[0], "line1") {
		t.Errorf("Expected first line to contain 'line1', got %q", lines[0])
	}

	// Second line should be indented and contain "line2"
	if !strings.Contains(lines[1], "line2") {
		t.Errorf("Expected second line to contain 'line2', got %q", lines[1])
	}

	// Third line should be indented and contain "line3"
	if !strings.Contains(lines[2], "line3") {
		t.Errorf("Expected third line to contain 'line3', got %q", lines[2])
	}
}
