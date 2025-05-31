package utils

import "strconv"

func StringToInt(s string, defaultValue int) int {
	if value, err := strconv.Atoi(s); err == nil {
		return value
	}
	return defaultValue
}

func StringToBool(s string, defaultValue bool) bool {
	if value, err := strconv.ParseBool(s); err == nil {
		return value
	}
	return defaultValue
}