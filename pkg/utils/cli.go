package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
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
