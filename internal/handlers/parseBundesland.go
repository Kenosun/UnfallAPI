package handlers

func parseBundesland(numstr string) string {
	switch numstr {
	case "01":
		return "Schleswig-Holstein"
	case "02":
		return "Hamburg"
	case "03":
		return "Niedersachsen"
	case "04":
		return "Bremen"
	case "05":
		return "Nordrhein-Westfalen"
	case "06":
		return "Hessen"
	case "07":
		return "Rheinland-Pfalz"
	case "08":
		return "Baden-Württemberg"
	case "09":
		return "Bayern"
	case "10":
		return "Saarland"
	case "11":
		return "Berlin"
	case "12":
		return "Brandenburg"
	case "13":
		return "Mecklenburg-Vorpommern"
	case "14":
		return "Sachsen"
	case "15":
		return "Sachsen-Anhalt"
	case "16":
		return "Thüringen"
	default:
		return ""
	}
}
