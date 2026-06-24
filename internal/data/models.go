package data

type UnfallStatistik struct {
	Unfallkategorie string `json:"unfallkategorie"`
	Ortslage        string `json:"ortslage"`
	Jahr            int    `json:"jahr"`
	Monat           int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl          int    `json:"anzahl"`
}

type UnfallStrassenverkehr struct {
	Strassenklasse string `json:"strassenklasse"`
	Ortslage       string `json:"ortslage"`
	Kategorie      string `json:"kategorie"`
	Jahr           int    `json:"jahr"`
	Monat          int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl         int    `json:"anzahl"`
}

type UnfallPersonenschaden struct {
	Unfalltyp   string `json:"unfalltyp"`
	Ortslage    string `json:"ortslage"`
	Schweregrad string `json:"schweregrad"`
	Kategorie   string `json:"kategorie"`
	Jahr        int    `json:"jahr"`
	Monat       int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl      int    `json:"anzahl"`
}

type UnfallVerunglueckte struct {
	Verkehrsart  string `json:"verkehrsart"`
	Ortslage     string `json:"ortslage"`
	Schweregrad  string `json:"schweregrad"`
	Geschlecht   string `json:"geschlecht"`
	Altersgruppe string `json:"altersgruppe"`
	Jahr         int    `json:"jahr"`
	Monat        int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl       int    `json:"anzahl"`
}

type UnfallFehlverhalten struct {
	Verkehrsart   string `json:"verkehrsart"`
	Fehlverhalten string `json:"fehlverhalten"`
	Jahr          int    `json:"jahr"`
	Monat         int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl        int    `json:"anzahl"`
}

type UnfallBeteiligung struct {
	Verkehrsart     string `json:"verkehrsart"`
	Unfallkategorie string `json:"unfallkategorie"`
	Ortslage        string `json:"ortslage"`
	Geschlecht      string `json:"geschlecht"`
	Altersgruppe    string `json:"altersgruppe"`
	Beteiligungsart string `json:"beteiligungsart"` // "Unfallbeteiligte" / "Hauptverursacher des Unfalls"
	Jahr            int    `json:"jahr"`
	Monat           int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl          int    `json:"anzahl"`
}

type UnfallStatistikBundesland struct {
	Bundesland      string `json:"bundesland"`
	Unfallkategorie string `json:"unfallkategorie"`
	Ortslage        string `json:"ortslage"`
	Jahr            int    `json:"jahr"`
	Monat           int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl          int    `json:"anzahl"`
}

type UnfallStrassenverkehrBundesland struct {
	Bundesland     string `json:"bundesland"`
	Strassenklasse string `json:"strassenklasse"`
	Ortslage       string `json:"ortslage"`
	Jahr           int    `json:"jahr"`
	Anzahl         int    `json:"anzahl"`
}

type UnfallVerunglueckteBundesland struct {
	Bundesland  string `json:"bundesland"`
	Ortslage    string `json:"ortslage"`
	Schweregrad string `json:"schweregrad"`
	Jahr        int    `json:"jahr"`
	Monat       int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl      int    `json:"anzahl"`
}

type Unfall struct {
	Bundesland                 string  `json:"bundesland"`
	Regierungsbezirk           string  `json:"regierungsbezirk"`
	Kreis                      string  `json:"kreis"`
	Gemeinde                   string  `json:"gemeinde"`
	Jahr                       int     `json:"jahr"`
	Monat                      int     `json:"monat"`
	Uhrzeit                    string  `json:"uhrzeit"`
	Wochentag                  string  `json:"wochentag"`
	Schweregrad                string  `json:"schweregrad"`
	Unfallart                  string  `json:"unfallart"`
	Unfalltyp                  string  `json:"unfalltyp"`
	Lichtverhaeltnis           string  `json:"lichtverhaeltnis"`
	MitFahrrad                 bool    `json:"mit_fahrrad"`
	MitPKW                     bool    `json:"mit_pkw"`
	MitFussgaenger             bool    `json:"mit_fussgaenger"`
	MitKraftrad                bool    `json:"mit_kraftrad"`
	MitGueterkraftfahrzeug     bool    `json:"mit_gueterkraftfahrzeug"`
	MitSonstigenVerkehrsmittel bool    `json:"mit_sonstigen_verkehrsmittel"`
	Strassenzustand            string  `json:"strassenzustand"`
	Latitude                   float64 `json:"latitude"`
	Longitude                  float64 `json:"longitude"`
}

type Ort struct {
	Bundesland          string  `json:"bundesland"`
	Regierungsbezirk    string  `json:"regierungsbezirk"`
	Kreis               string  `json:"kreis"`
	Gemeinde            string  `json:"gemeinde"`
	Name                string  `json:"name"`
	Gemeindeverband     string  `json:"gemeindeverband"`
	Landkreis           string  `json:"landkreis"`
	Postleitzahl        string  `json:"postleitzahl"`
	Flaeche             float64 `json:"flaeche"`
	Bevoelkerung        int     `json:"bevoelkerung"`
	Maennlich           int     `json:"maennlich"`
	Weiblich            int     `json:"weiblich"`
	Reisegebiet         string  `json:"reisegebiet"`
	Verstaedterungsgrad string  `json:"verstaedterungsgrad"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
}
