package helper

func ParseUnfalltyp(num int) string {
	switch num {
	case 1:
		return "Fahrunfall"
	case 2:
		return "Abbiegeunfall"
	case 3:
		return "Einbiegen/Kreuzen-Unfall"
	case 4:
		return "Überschreiten-Unfall"
	case 5:
		return "Unfall durch ruhenden Verkehr"
	case 6:
		return "Unfall im Längsverkehr"
	case 7:
		return "sonstiger Unfall"
	default:
		return ""
	}
}
