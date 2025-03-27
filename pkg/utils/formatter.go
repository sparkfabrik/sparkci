package utils

import (
	"fmt"
	"strings"
	"time"
)

// For testing purposes
var timeNow = time.Now

// GitLabSectionStart begins a collapsible section in GitLab CI output
func GitLabSectionStart(sectionTitle string, sectionDescription string) {
	if sectionDescription == "" {
		sectionDescription = sectionTitle
	}

	timestamp := timeNow().Unix()
	fmt.Printf("section_start:%d:%s[collapsed=true]\r\033[0K%s\n",
		timestamp, sectionTitle, sectionDescription)
}

// GitLabSectionEnd closes a section in GitLab CI output
func GitLabSectionEnd(sectionTitle string) {
	timestamp := timeNow().Unix()
	fmt.Printf("section_end:%d:%s\r\033[0K\n", timestamp, sectionTitle)
}

func GitLabPrintBanner(text string) {
	if text != "" {
		border := "+" + strings.Repeat("-", len(text)+2) + "+"
		fmt.Println(border)
		fmt.Println("| " + text + " |")
		fmt.Println(border)
	}
}
