package utils

func GetOrDefault(str string, defaultStr string) string {
	if str != "" {
		return str
	}
	return defaultStr
}

func IsAnyBlank(strList ...string) bool {
	if strList == nil || len(strList) <= 0 {
		return false
	}
	for _, str := range strList {
		if str == "" {
			return true
		}
	}
	return false
}
