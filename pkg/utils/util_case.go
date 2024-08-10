package utils

func CaseToPoint[T any](ele T) *T {
	return &ele
}

func Ptr[T any](t T) *T {
	return CaseToPoint(t)
}
