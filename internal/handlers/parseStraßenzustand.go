package handlers

func parseStraßenzustand(num int) string {
	switch num {
	case 0:
		return "trocken"
	case 1:
		return "nass/feucht/schlüpfrig"
	case 2:
		return "winterglatt"
	default:
		return ""
	}
}
