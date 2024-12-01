package utils

import "golang.org/x/exp/constraints"

// NOTE: I'd like to use just simple `b - a` for descending sort but the
// problem will rise with casting MAX_UINT_64 to INT_32 that CAN lead to
// potential (0xFFFFFFFF) value as return from `MAX_UINT_64 - 0` and I just
// want to be sure that this part won't cause any trouble in the future
func DescendingSort[T constraints.Ordered](a, b T) int {
	if a > b {
		return -1
	} else if a < b {
		return 1
	} else {
		return 0
	}
}

func AscendingSort[T constraints.Ordered](a, b T) int {
	if b > a {
		return -1
	} else if b < a {
		return 1
	} else {
		return 0
	}
}
