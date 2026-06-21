package helper

func ParseSchweregrad(num int) string {
	switch num {
	case 1:
		return "Unfall mit Getöteten"
	case 2:
		return "Unfall mit Schwerverletzten"
	case 3:
		return "Unfall mit Leichtverletzten"
	default:
		return ""
	}
}
