package utils

func IfNotNilPanic(errs ...error) {
	if errs == nil || len(errs) <= 0 {
		return
	}
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
