package utils

func singleToSlice(s string) []string {
	return []string{s}
}

func ToSlice[T any](v T) []T {
	return []T{v}
}
