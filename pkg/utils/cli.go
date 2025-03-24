package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ANSI escape codes for text formatting
const (
	AnsiReset  = "\033[0m"
	AnsiBold   = "\033[1m"
	AnsiRed    = "\033[31m"
	AnsiGreen  = "\033[32m"
	AnsiYellow = "\033[33m"
	AnsiBlue   = "\033[34m"
)

// PrintKeyValue prints key-value pairs with padding for alignment.
// The pad parameter specifies additional spaces after the key and colon.
func PrintKeyValue(w io.Writer, pad int, key, value string) {
	if pad < 0 {
		pad = 0
	}
	fmt.Fprintf(w, "%s: %s%s\n", key, strings.Repeat(" ", pad), value)
}

// PrintMap prints all key-value pairs in a map with consistent alignment
func PrintMap(m map[string]string) {
	// First determine the longest key length
	maxKeyLen := 0
	for k := range m {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	// Calculate padding for each key based on max length
	for k, v := range m {
		pad := (maxKeyLen - len(k))
		PrintKeyValue(os.Stdout, pad, k, v)
	}
}

// PrintStyledKeyValue prints a key-value pair with custom padding and styling
func PrintStyledKeyValue(w io.Writer, padLen int, key string, value string, valueBold bool) {
	// Format the key with trailing colon and padding
	paddedKey := fmt.Sprintf("%-*s", padLen, key+":")

	// Apply bold formatting to value if requested
	formattedValue := value
	if valueBold {
		formattedValue = AnsiBold + value + AnsiReset
	}

	fmt.Fprintf(w, "%s %s\n", paddedKey, formattedValue)
}

// PrintSectionHeader prints a section header with bold formatting
func PrintSectionHeader(w io.Writer, title string) {
	fmt.Fprintf(w, "%s%s%s\n", AnsiBold, title, AnsiReset)
}

func PrintFormattedVars(title string, vars map[string]string) {
	// Auto adjust the padLen based on the longest key.
	padLen := 0
	for k := range vars {
		if len(k) > padLen {
			padLen = len(k) + 2 // Add space for colon and padding
		}
	}

	if title != "" {
		PrintSectionHeader(os.Stdout, title)
	}

	// Print each variable with consistent formatting
	for k, v := range vars {
		PrintStyledKeyValue(os.Stdout, padLen, k, v, true)
	}
}

func PrintVarGroup(title string, vars map[string]string) {
	PrintFormattedVars(title, vars)
}
