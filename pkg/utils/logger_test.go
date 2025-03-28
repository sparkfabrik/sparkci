package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerOutput(t *testing.T) {
	// Use a buffer to capture output
	var buf bytes.Buffer

	// Create a test logger that writes to our buffer
	logger := NewLogger(
		WithOutput(&buf),
		WithTimeFormat(""), // Disable timestamp for testing
		WithPrefix("test"),
		WithLevel(LogLevelDebug),
	)

	// Test cases
	testCases := []struct {
		name     string
		logFunc  func(string, ...interface{})
		message  string
		expected string
	}{
		{
			name:     "Debug log",
			logFunc:  logger.Debug,
			message:  "debug message",
			expected: "debug message\n",
		},
		{
			name:     "Info log",
			logFunc:  logger.Info,
			message:  "info message",
			expected: "info message\n",
		},
		{
			name:     "Warn log",
			logFunc:  logger.Warn,
			message:  "warning message",
			expected: "warning message\n",
		},
		{
			name:     "Error log",
			logFunc:  logger.Error,
			message:  "error message",
			expected: "error message\n",
		},
		{
			name:     "Formatted message",
			logFunc:  logger.Info,
			message:  "formatted %s %d",
			expected: "formatted string 42\n",
		},
		{
			name:     "Multiline message",
			logFunc:  logger.Info,
			message:  "line1\nline2\nline3",
			expected: "line1\nline2\nline3\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear buffer before each test
			buf.Reset()

			// Call the log function
			if tc.name == "Formatted message" {
				tc.logFunc(tc.message, "string", 42)
			} else {
				tc.logFunc(tc.message)
			}

			// Check output
			got := buf.String()
			if got != tc.expected {
				t.Errorf("Expected output %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestLogLevel(t *testing.T) {
	var buf bytes.Buffer

	// Create logger with INFO level
	logger := NewLogger(
		WithOutput(&buf),
		WithLevel(LogLevelInfo),
	)

	// Debug messages should be ignored
	logger.Debug("This debug message should not appear")
	if buf.Len() > 0 {
		t.Error("Debug message was logged when level is INFO")
	}

	// Info messages should be logged
	buf.Reset()
	logger.Info("This info message should appear")
	if buf.Len() == 0 {
		t.Error("Info message was not logged when level is INFO")
	}

	// Verify message contains only the content without prefixes
	output := buf.String()
	if !strings.Contains(output, "This info message should appear") {
		t.Errorf("Info message not found in output: %q", output)
	}

	// Make sure there's no prefix in the output
	if strings.Contains(output, "sparkci") ||
		strings.Contains(output, "INFO") ||
		strings.Contains(output, "[") {
		t.Errorf("Output should not contain prefix or log level indicators: %q", output)
	}
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(
		WithOutput(&buf),
	)

	// Add fields to logger
	fields := map[string]interface{}{
		"user": "testuser",
		"id":   123,
	}

	fieldLogger := logger.WithFields(fields)

	// Log with fields
	fieldLogger.Info("Test message with fields")

	// Check output - should only contain the message without field info
	output := buf.String()
	if output != "Test message with fields\n" {
		t.Errorf("Expected simple output without fields, got: %q", output)
	}
}

// Avoid testing Fatal since it calls os.Exit
func TestFatalWithoutExit(t *testing.T) {
	// Override os.Exit temporarily to avoid test termination
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	exitCalled := false
	osExit = func(code int) {
		exitCalled = true
	}

	var buf bytes.Buffer
	logger := NewLogger(WithOutput(&buf))

	logger.Fatal("Fatal error message")

	if !exitCalled {
		t.Error("os.Exit was not called by Fatal")
	}

	output := buf.String()
	if output != "Fatal error message\n" {
		t.Errorf("Expected simple output for fatal error, got: %q", output)
	}
}
