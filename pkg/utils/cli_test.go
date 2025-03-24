package utils

import (
	"bytes"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintKeyValue(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		pad      int
		expected string
	}{
		{
			name:     "No padding needed",
			key:      "Key",
			value:    "Value",
			pad:      0,
			expected: "Key: Value\n",
		},
		{
			name:     "With padding",
			key:      "Key",
			value:    "Value",
			pad:      5,
			expected: "Key:      Value\n",
		},
		{
			name:     "Long key with padding",
			key:      "VeryLongKey",
			value:    "Value",
			pad:      2,
			expected: "VeryLongKey:   Value\n",
		},
		{
			name:     "Empty key",
			key:      "",
			value:    "Value only",
			pad:      3,
			expected: ":    Value only\n",
		},
		{
			name:     "Empty value",
			key:      "Key only",
			value:    "",
			pad:      1,
			expected: "Key only:  \n",
		},
		{
			name:     "Negative padding (should be treated as 0)",
			key:      "Key",
			value:    "Value",
			pad:      -1,
			expected: "Key: Value\n",
		},
		{
			name:     "Large padding",
			key:      "Key",
			value:    "Value",
			pad:      10,
			expected: "Key:           Value\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture output
			buffer := bytes.Buffer{}

			// Call the function with our buffer
			PrintKeyValue(&buffer, tt.pad, tt.key, tt.value)

			// Verify the output matches expected
			assert.Equal(t, tt.expected, buffer.String())
		})
	}
}

func TestPrintMap(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[string]string
		verifyFn func(t *testing.T, output string)
	}{
		{
			name:     "Empty map",
			inputMap: map[string]string{},
			verifyFn: func(t *testing.T, output string) {
				assert.Equal(t, "", strings.TrimSpace(output))
			},
		},
		{
			name: "One key-value pair",
			inputMap: map[string]string{
				"Key": "Value",
			},
			verifyFn: func(t *testing.T, output string) {
				assert.Contains(t, output, "Key: Value")
			},
		},
		{
			name: "Multiple key-value pairs of same length",
			inputMap: map[string]string{
				"Key1": "Value1",
				"Key2": "Value2",
				"Key3": "Value3",
			},
			verifyFn: func(t *testing.T, output string) {
				lines := strings.Split(strings.TrimSpace(output), "\n")
				assert.Equal(t, 3, len(lines))

				// All keys should have the same padding since they're the same length
				for _, line := range lines {
					assert.True(t, strings.Contains(line, "Key"), "Line should contain 'Key'")
					assert.True(t, strings.Contains(line, ": "), "Line should contain ': '")
				}

				// Verify content is present
				assert.Contains(t, output, "Key1: Value1")
				assert.Contains(t, output, "Key2: Value2")
				assert.Contains(t, output, "Key3: Value3")
			},
		},
		{
			name: "Keys of different lengths",
			inputMap: map[string]string{
				"A":           "ValueA",
				"LongerKey":   "ValueB",
				"VeryLongKey": "ValueC",
			},
			verifyFn: func(t *testing.T, output string) {
				lines := strings.Split(strings.TrimSpace(output), "\n")
				assert.Equal(t, 3, len(lines))

				// Sort lines to make testing deterministic
				sort.Strings(lines)

				// Each line should have the format "Key: [padding]Value"
				for _, line := range lines {
					parts := strings.SplitN(line, ": ", 2)
					assert.Equal(t, 2, len(parts), "Line should have key: value format")
				}

				// The shorter keys should have more padding after the colon
				aLine := lines[0] // "A: ..." (after sorting)
				assert.True(t, strings.HasPrefix(aLine, "A: "))

				// Check that values align
				// We should find the position of each value in its line
				aValuePos := strings.Index(aLine, "ValueA")

				for i := 1; i < len(lines); i++ {
					line := lines[i]
					if strings.Contains(line, "ValueB") {
						bValuePos := strings.Index(line, "ValueB")
						assert.Equal(t, aValuePos, bValuePos, "Values should be aligned")
					}
					if strings.Contains(line, "ValueC") {
						cValuePos := strings.Index(line, "ValueC")
						assert.Equal(t, aValuePos, cValuePos, "Values should be aligned")
					}
				}
			},
		},
		{
			name: "With empty values",
			inputMap: map[string]string{
				"Key1": "",
				"Key2": "Value2",
			},
			verifyFn: func(t *testing.T, output string) {
				assert.Contains(t, output, "Key1: ")
				assert.Contains(t, output, "Key2: Value2")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call PrintMap
			PrintMap(tt.inputMap)

			// Restore stdout and read captured output
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Run test-specific verifications
			tt.verifyFn(t, output)
		})
	}
}

// TestPrintMapWithMock tests PrintMap with a mock io.Writer
func TestPrintMapWithMock(t *testing.T) {
	// Create a custom version of PrintMap that uses the provided writer
	customPrintMap := func(w io.Writer, m map[string]string) {
		maxKeyLen := 0
		for k := range m {
			if len(k) > maxKeyLen {
				maxKeyLen = len(k)
			}
		}

		for k, v := range m {
			pad := maxKeyLen - len(k)
			PrintKeyValue(w, pad, k, v)
		}
	}

	testMap := map[string]string{
		"Key1":      "Value1",
		"LongerKey": "Value2",
	}

	var buf bytes.Buffer
	customPrintMap(&buf, testMap)
	output := buf.String()

	// Verify alignment
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Equal(t, 2, len(lines))

	// Verify content
	assert.Contains(t, output, "Key1:")
	assert.Contains(t, output, "Value1")
	assert.Contains(t, output, "LongerKey:")
	assert.Contains(t, output, "Value2")

	// The shorter key should have more padding
	key1Line := ""
	longerKeyLine := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "Key1:") {
			key1Line = line
		} else if strings.HasPrefix(line, "LongerKey:") {
			longerKeyLine = line
		}
	}

	// Extract the space after the colon
	key1Spaces := strings.TrimPrefix(key1Line, "Key1:")
	key1Value := strings.TrimSpace(key1Spaces)
	longerKeySpaces := strings.TrimPrefix(longerKeyLine, "LongerKey:")
	longerKeyValue := strings.TrimSpace(longerKeySpaces)

	// Key1 should have more padding spaces than LongerKey
	assert.Greater(t, len(key1Spaces), len(longerKeySpaces))

	// But values should be intact
	assert.Equal(t, "Value1", key1Value)
	assert.Equal(t, "Value2", longerKeyValue)
}
func TestPrintFormattedVars(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		vars     map[string]string
		expected []string
	}{
		{
			name:  "No title, single variable",
			title: "",
			vars: map[string]string{
				"Key": "Value",
			},
			expected: []string{
				"Key:  \x1b[1mValue",
			},
		},
		{
			name:  "With title, single variable",
			title: "Header",
			vars: map[string]string{
				"Key": "Value",
			},
			expected: []string{
				"\033[1mHeader\033[0m",
				"Key:  \x1b[1mValue",
			},
		},
		{
			name:  "Multiple variables",
			title: "Variables",
			vars: map[string]string{
				"Key1": "Value1",
				"Key2": "Value2",
			},
			expected: []string{
				"\033[1mVariables\033[0m",
				"Key1:  \x1b[1mValue1",
				"Key2:  \x1b[1mValue2",
			},
		},
		{
			name:  "Empty variables",
			title: "Empty",
			vars:  map[string]string{},
			expected: []string{
				"\033[1mEmpty\033[0m",
			},
		},
		{
			name:  "No title, multiple variables",
			title: "",
			vars: map[string]string{
				"Short":     "Value1",
				"LongerKey": "Value2",
			},
			expected: []string{
				"Short:      \x1b[1mValue1\x1b[0m\n",
				"LongerKey:  \x1b[1mValue2\x1b[0m\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call PrintFormattedVars
			PrintFormattedVars(tt.title, tt.vars)

			// Restore stdout and read captured output
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify output
			for _, expectedLine := range tt.expected {
				assert.Contains(t, output, expectedLine)
			}
		})
	}
}

func TestPrintVarGroup(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		vars     map[string]string
		expected []string
	}{
		{
			name:  "Group with title",
			title: "GroupTitle",
			vars: map[string]string{
				"Key1": "Value1",
				"Key2": "Value2",
			},
			expected: []string{
				"\033[1mGroupTitle\033[0m",
				"Key1:  \x1b[1mValue1",
				"Key2:  \x1b[1mValue2",
			},
		},
		{
			name:  "Group without title",
			title: "",
			vars: map[string]string{
				"Key1": "Value1",
				"Key2": "Value2",
			},
			expected: []string{
				"Key1:  \x1b[1mValue1",
				"Key2:  \x1b[1mValue2",
			},
		},
		{
			name:  "Empty group",
			title: "EmptyGroup",
			vars:  map[string]string{},
			expected: []string{
				"\033[1mEmptyGroup\033[0m",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call PrintVarGroup
			PrintVarGroup(tt.title, tt.vars)

			// Restore stdout and read captured output
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify output
			for _, expectedLine := range tt.expected {
				assert.Contains(t, output, expectedLine)
			}
		})
	}
}
