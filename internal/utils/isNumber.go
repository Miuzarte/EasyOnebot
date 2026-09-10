package utils

var numberTestMap = [256]bool{
	'+': true, '-': true,
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
}

func IsNumber(s string, additional ...byte) bool {
	if len(s) == 0 {
		return false
	}

	if len(additional) == 0 {
		for i := range len(s) {
			if !numberTestMap[s[i]] {
				return false
			}
		}
	} else {
		testMap := make([]bool, 256)
		copy(testMap, numberTestMap[:])
		for i := range additional {
			testMap[additional[i]] = true
		}
		for i := range len(s) {
			if !testMap[s[i]] {
				return false
			}
		}
	}

	return true
}
