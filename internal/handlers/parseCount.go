package handlers

import (
	"strconv"
	"strings"
)

func parseCount(valStr string) (int, bool) {
	valStr = strings.TrimSpace(valStr)
	if valStr == "" || valStr == "-" {
		return 0, false
	}

	// try parsing directly
	if count, err := strconv.Atoi(valStr); err == nil {
		return count, true
	}

	// strip thousand separator (27.348 -> 27348)
	cleaned := strings.ReplaceAll(valStr, ".", "")
	if count, err := strconv.Atoi(cleaned); err == nil {
		return count, true
	}

	return -1, false
}
