package convert

import "strconv"

func ParseInt(s string, defaultValue int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return val
}
