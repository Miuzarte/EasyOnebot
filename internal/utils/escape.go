package utils

import (
	"strings"
)

var (
	htmlEscape = strings.NewReplacer(
		"&", "&amp;",
		",", "&#44;",
		"[", "&#91;",
		"]", "&#93;",
	)
	htmlUnescape = strings.NewReplacer(
		"&amp;", "&",
		"&#44;", ",",
		"&#91;", "[",
		"&#93;", "]",
	)
)

func EscapeCqCode(str string) string {
	return htmlEscape.Replace(str)
}

func UnescapeCqCode(str string) string {
	return htmlUnescape.Replace(str)
}
