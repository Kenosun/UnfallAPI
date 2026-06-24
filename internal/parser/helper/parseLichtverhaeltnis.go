package helper

func ParseLichtverhaeltnis(num int) string {
	switch num {
	case 0:
		return "Tageslicht"
	case 1:
		return "Dämmerung"
	case 2:
		return "Dunkelheit"
	default:
		return ""
	}
}
