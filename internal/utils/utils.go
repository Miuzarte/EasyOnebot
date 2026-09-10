package utils

import (
	"fmt"

	"github.com/jinzhu/copier"
)

// AnyCopy 复制 any 到 any
func AnyCopy[T any](from any) (to *T) {
	to = new(T)
	err := copier.Copy(to, from)
	if err != nil {
		panic(fmt.Sprintf("AnyCopy: failed to copy from %T to %T: %v", from, to, err))
	}
	return to
}
