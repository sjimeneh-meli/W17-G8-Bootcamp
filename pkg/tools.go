package tools

import (
	"fmt"
	"strings"
)

func SliceToString[T any](slice []T, delim string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(slice)), delim), "[]")
}

func FillNewSlice[T any](size int, value T) []T {
	sliceFilled := make([]T, size)
	for i := range size {
		sliceFilled[i] = value
	}
	return sliceFilled
}

func GetSlicesOfKeyAndValuesFromMap(m map[any]any) ([]any, []any) {
	keys := make([]any, 0, len(m))
	values := make([]any, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
