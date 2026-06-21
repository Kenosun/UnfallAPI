package helper

func ParseUnfallart(num int) string {
	switch num {
	case 1:
		return "Zusammenstoß mit anfahrendem/anhaltendem/ruhendem Fahrzeug"
	case 2:
		return "Zusammenstoß mit vorausfahrendem/wartendem Fahrzeug"
	case 3:
		return "Zusammenstoß mit seitlich in gleicher Richtung fahrendem Fahrzeug"
	case 4:
		return "Zusammenstoß mit entgegenkommendem Fahrzeug"
	case 5:
		return "Zusammenstoß mit einbiegendem/kreuzendem Fahrzeug"
	case 6:
		return "Zusammenstoß zwischen Fahrzeug und Fußgänger"
	case 7:
		return "Aufprall auf Fahrbahnhindernis"
	case 8:
		return "Abkommen von Fahrbahn nach rechts"
	case 9:
		return "Abkommen von Fahrbahn nach links"
	case 0:
		return "Unfall anderer Art"
	default:
		return ""
	}
}
