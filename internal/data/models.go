package data

type UnfallStatistik struct {
	Kategorie string `json:"kategorie"`
	Ortslage  string `json:"ortslage"`
	Jahr      int    `json:"jahr"`
	Monat     int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl    int    `json:"anzahl"`
}

type UnfallStraßenverkehr struct {
	Straßenklasse string `json:"straßenklasse"`
	Ortslage      string `json:"ortslage"`
	Kategorie     string `json:"kategorie"`
	Jahr          int    `json:"jahr"`
	Monat         int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl        int    `json:"anzahl"`
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

type UnfallVerunglückte struct {
	Verkehrsart  string `json:"verkehrsart"`
	Ortslage     string `json:"ortslage"`
	Kategorie    string `json:"kategorie"`
	Geschlecht   string `json:"geschlecht"`
	Altersgruppe string `json:"altersgruppe"`
	Jahr         int    `json:"jahr"`
	Monat        int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl       int    `json:"anzahl"`
}

type UnfallFehlverhalten struct {
	Verkehrsart string `json:"verkehrsart"`
	Kategorie   string `json:"kategorie"`
	Jahr        int    `json:"jahr"`
	Monat       int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl      int    `json:"anzahl"`
}

type UnfallBeteiligung struct {
	Verkehrsart     string `json:"verkehrsart"`
	Kategorie       string `json:"kategorie"`
	Ortslage        string `json:"ortslage"`
	Geschlecht      string `json:"geschlecht"`
	Altersgruppe    string `json:"altersgruppe"`
	Beteiligungsart string `json:"beteiligungsart"` // "Unfallbeteiligte" / "Hauptverursacher des Unfalls"
	Jahr            int    `json:"jahr"`
	Monat           int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl          int    `json:"anzahl"`
}

type UnfallStatistikBundesland struct {
	Bundesland string `json:"bundesland"`
	Kategorie  string `json:"kategorie"`
	Ortslage   string `json:"ortslage"`
	Jahr       int    `json:"jahr"`
	Monat      int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl     int    `json:"anzahl"`
}

type UnfallStraßenverkehrBundesland struct {
	Bundesland    string `json:"bundesland"`
	Straßenklasse string `json:"straßenklasse"`
	Ortslage      string `json:"ortslage"`
	Jahr          int    `json:"jahr"`
	Monat         int    `json:"monat"` // 1-12 for months, 0 for full year data
	Anzahl        int    `json:"anzahl"`
}

type UnfallVerunglückteBundesland struct {
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
	Lichtverhältnis            string  `json:"lichtverhältnis"`
	MitFahrrad                 bool    `json:"mit_fahrrad"`
	MitPKW                     bool    `json:"mit_pkw"`
	MitFußgänger               bool    `json:"mit_fußgänger"`
	MitKraftrad                bool    `json:"mit_kraftrad"`
	MitGüterkraftfahrzeug      bool    `json:"mit_güterkraftfahrzeug"`
	MitSonstigenVerkehrsmittel bool    `json:"mit_sonstigen_verkehrsmittel"`
	IstStraße                  bool    `json:"ist_straße"`
	Straßenzustand             string  `json:"straßenzustand"`
	Latitude                   float64 `json:"latitude"`
	Longitude                  float64 `json:"longitude"`
}

type Ort struct {
	Bundesland         string  `json:"bundesland"`
	Regierungsbezirk   string  `json:"regierungsbezirk"`
	Kreis              string  `json:"kreis"`
	Gemeinde           string  `json:"gemeinde"`
	Name               string  `json:"name"`
	Gemeindeverband    string  `json:"gemeindeverband"`
	Landkreis          string  `json:"landkreis"`
	Postleitzahl       string  `json:"postleitzahl"`
	Fläche             float64 `json:"fläche"`
	Bevölkerung        int     `json:"bevölkerung"`
	Männlich           int     `json:"männlich"`
	Weiblich           int     `json:"weiblich"`
	Reisegebiet        string  `json:"reisegebiet"`
	Verstädterungsgrad string  `json:"verstädterungsgrad"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
}
