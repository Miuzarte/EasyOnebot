package utils

import "strings"

func StringsToLowerMulti(s []string) []string {
	for i := range s {
		s[i] = strings.ToLower(s[i])
	}
	return s
}

func StringsContainsMulti(s string, substrs []string) bool {
	s = strings.ToLower(s)
	substrs = StringsToLowerMulti(substrs)
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
