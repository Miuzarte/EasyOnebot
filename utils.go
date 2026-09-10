package EasyOnebot

import "strconv"

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr
}

type floatPoint interface {
	~float32 | ~float64
}

type number interface {
	integer | floatPoint
}

func itoa[T number](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func utoa[T number](i T) string {
	return strconv.FormatUint(uint64(i), 10)
}
