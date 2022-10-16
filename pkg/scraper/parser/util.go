package parser

import (
	"strings"
)

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

func findKeyByValueInMultiMap[K, V comparable](m map[K][]V, val V) (K, bool) {
	for k, vs := range m {
		for _, v := range vs {
			if val == v {
				return k, true
			}
		}
	}

	var k K
	return k, false
}
