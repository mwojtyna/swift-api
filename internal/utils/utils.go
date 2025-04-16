package utils

func Map[T, V any](slice []T, fn func(T) V) []V {
	result := make([]V, len(slice))
	for i, t := range slice {
		result[i] = fn(t)
	}
	return result
}

const hqPartLen = 8

// Returns whether the bank is the headquarters and if not, also returns the bank's headquarters code assuming they exist
func IsSwiftCodeHq(code string) (bool, string) {
	if code[hqPartLen:] == "XXX" {
		return true, ""
	} else {
		return false, code[:hqPartLen] + "XXX"
	}
}
