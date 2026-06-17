package handlers

type HeaderYearMonth struct {
	Year  int
	Month int
}

type HeaderGenderAge struct {
	Geschlecht   string
	Altersgruppe string
}

type HeaderUnfallBeteiligung struct {
	Geschlecht      string
	Altersgruppe    string
	Beteiligungsart string
}
