package utils

func RuneCut(s string, n ...int) string {
	const defaultN = 4
	if len(n) == 0 {
		if len([]rune(s)) > defaultN {
			return string([]rune(s)[:defaultN]) + "..."
		}
	} else if len([]rune(s)) > n[0] {
		return string([]rune(s)[:n[0]]) + "..."
	}
	return s
}
