package helper

import "strings"

func ParseMonthToInt(month string) int {
	switch strings.ToLower(month) {
	case "januar":
		return 1
	case "februar":
		return 2
	case "märz":
		return 3
	case "april":
		return 4
	case "mai":
		return 5
	case "juni":
		return 6
	case "juli":
		return 7
	case "august":
		return 8
	case "september":
		return 9
	case "oktober":
		return 10
	case "november":
		return 11
	case "dezember":
		return 12
	default:
		return -1
	}
}
