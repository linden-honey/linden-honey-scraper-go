package parser

import (
	"strings"
)

// substringAfterLast util function to get last substring after some inclusion
func substringAfterLast(s, substr string) string {
	if len(s) == 0 {
		return s
	}

	startIndex := strings.LastIndex(s, substr)
	if len(substr) == 0 || startIndex == -1 {
		return ""
	}

	return s[startIndex+len(substr):]
}
