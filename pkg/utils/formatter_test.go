package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	// Redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function that prints to stdout
	f()

	// Reset stdout and get the output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestGitLabSectionStart(t *testing.T) {
	// Mock the time.Now() function to return a fixed timestamp
	now := time.Now()
	expectedTimestamp := now.Unix()

	// Store the original time.Now function and restore it after the test
	oldTimeNow := timeNow
	timeNow = func() time.Time { return now }
	defer func() { timeNow = oldTimeNow }()

	t.Run("With Title and Description", func(t *testing.T) {
		title := "test-section"
		desc := "Test Section Description"

		output := captureOutput(func() {
			GitLabSectionStart(title, desc)
		})

		expectedOutput := fmt.Sprintf("section_start:%d:%s[collapsed=true]\r\033[0K%s\n",
			expectedTimestamp, title, desc)
		assert.Equal(t, expectedOutput, output)
	})

	t.Run("With Title Only", func(t *testing.T) {
		title := "test-section"

		output := captureOutput(func() {
			GitLabSectionStart(title, "")
		})

		expectedOutput := fmt.Sprintf("section_start:%d:%s[collapsed=true]\r\033[0K%s\n",
			expectedTimestamp, title, title)
		assert.Equal(t, expectedOutput, output)
	})

	t.Run("Format Contains Required Elements", func(t *testing.T) {
		title := "test-section"
		desc := "Test Section Description"

		output := captureOutput(func() {
			GitLabSectionStart(title, desc)
		})

		// Check that the output contains the necessary elements
		assert.True(t, strings.Contains(output, "section_start:"))
		assert.True(t, strings.Contains(output, title))
		assert.True(t, strings.Contains(output, desc))
		assert.True(t, strings.Contains(output, "[collapsed=true]"))
		assert.True(t, strings.Contains(output, "\r\033[0K")) // Terminal control sequence

		// Check the format follows GitLab CI spec
		parts := strings.SplitN(output, ":", 3)
		assert.Equal(t, "section_start", parts[0])
		assert.Contains(t, parts[1], fmt.Sprintf("%d", expectedTimestamp))
		assert.Contains(t, parts[2], title+"[collapsed=true]")
	})
}

func TestGitLabPrintBanner(t *testing.T) {
	t.Run("With Text", func(t *testing.T) {
		text := "Hello World"

		output := captureOutput(func() {
			GitLabPrintBanner(text)
		})

		border := "+" + strings.Repeat("-", len(text)+2) + "+"
		expected := border + "\n| " + text + " |\n" + border + "\n"
		assert.Equal(t, expected, output)
	})

	t.Run("With Empty Text", func(t *testing.T) {
		output := captureOutput(func() {
			GitLabPrintBanner("")
		})

		// Should not print anything if text is empty
		assert.Equal(t, "", output)
	})

	t.Run("Banner Format", func(t *testing.T) {
		text := "Test Banner"

		output := captureOutput(func() {
			GitLabPrintBanner(text)
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")

		// Verify we have exactly 3 lines
		assert.Equal(t, 3, len(lines), "Banner should have 3 lines")

		// First and last line should be identical borders with correct length
		expectedBorder := "+" + strings.Repeat("-", len(text)+2) + "+"
		assert.Equal(t, expectedBorder, lines[0], "First line should be a border")
		assert.Equal(t, expectedBorder, lines[2], "Last line should be a border")

		// Middle line should contain the text with proper formatting
		expectedMiddle := "| " + text + " |"
		assert.Equal(t, expectedMiddle, lines[1], "Middle line should contain the text with proper padding")
	})
}
