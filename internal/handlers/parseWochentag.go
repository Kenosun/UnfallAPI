package handlers

func parseWochentag(num int) string {
	switch num {
	case 1:
		return "Sonntag"
	case 2:
		return "Montag"
	case 3:
		return "Dienstag"
	case 4:
		return "Mittwoch"
	case 5:
		return "Donnerstag"
	case 6:
		return "Freitag"
	case 7:
		return "Samstag"
	default:
		return ""
	}
}
