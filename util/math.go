package util

func Gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func Gcd64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
