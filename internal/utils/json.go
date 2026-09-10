package utils

import (
	"encoding/json"
	"fmt"
)

func ToJson(a any) string {
	j, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprint(a)
	}
	return string(j)
}
