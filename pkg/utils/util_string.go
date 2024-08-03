package utils

func GetOrDefault(str string, defaultStr string) string {
	if str != "" {
		return str
	}
	return defaultStr
}
