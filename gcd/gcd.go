/*
Package gcd privides method to check if enough sample size for percentile
is provided, implementing Euclidean algorithm
*/
package gcd

// GCD takes two uint64 and returns their greatest common divisor as uint64
func GCD(a, b uint64) uint64 {
	for b != 0 {
		r := a % b
		a = b
		b = r
	}
	return a
}

func swap(a, b *uint64) {
	tmp := *a
	*a = *b
	*b = tmp
}

// swap is not necessary for Euclidean algorithm,
// because r is never larger then b
func swapIfSamller(a, b *uint64) {
	if *a < *b {
		swap(a, b)
	}
}
