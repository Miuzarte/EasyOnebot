package utils

import (
	"encoding/base64"
	"fmt"
)

func MarshalString(str any) string {
	switch s := str.(type) {
	case string:
		return s
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	case []byte:
		return string(s)
	default:
		return fmt.Sprint(str)
	}
}

// 只能在构造时使用, (不能在发送时使用)?
func MarshalFile(file any) (s string) {
	const base64Prefix = "base64://"
	switch file := file.(type) {
	case string:
		s = file
	case []byte: // data
		b64Data := make([]byte, len(base64Prefix)+base64.StdEncoding.EncodedLen(len(file)))
		copy(b64Data[:len(base64Prefix)], base64Prefix)
		base64.StdEncoding.Encode(b64Data[len(base64Prefix):], file)
		return string(b64Data)
	default:
		panic("unsupported file type, see [message.FileT]") // [message.FileT]
	}
	return
}
