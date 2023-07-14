package gtools

import "golang.org/x/exp/constraints"

func Ternary[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}

func LimitValue[T constraints.Ordered](v T, min T, max T) T {
	if v < min {
		return v
	}
	if v > max {
		return max
	}

	return v
}
