package slicez

import "github.com/zev-zakaryan/go-util/conv"

func Clone[T any](a []T) []T {
	return append([]T(nil), a...)
}

// Copy cast any slice to target slice.
func CloneCast[T any, U any](a []U) []T {
	b := make([]T, len(a))
	for i := range a {
		b[i] = conv.ToForce[T](a[i])
	}
	return b
}
