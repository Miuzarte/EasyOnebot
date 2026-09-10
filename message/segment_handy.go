package message

import (
	"fmt"
)

/*
Textf format 纯文本

orig: [Text]
*/
func Textf(format string, a ...any) Segment {
	return Text(fmt.Sprintf(format, a...))
}
