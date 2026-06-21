package handlers

import (
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Ort struct {
	Bundesland         string
	Regierungsbezirk   string
	Kreis              string
	Gemeinde           string
	Name               string
	Gemeindeverband    string
	Landkreis          string
	Postleitzahl       string
	Fläche             float64
	Bevölkerung        int
	Männlich           int
	Weiblich           int
	Reisegebiet        string
	Verstädterungsgrad string
	Latitude           float64
	Longitude          float64
}

func normalizeValue(val string) string {
	val = strings.ReplaceAll(val, " ", "")      // remove spaces
	val = strings.ReplaceAll(val, "\u00a0", "") // handle non-breaking spaces
	return val
}

func ParseOrt() ([]Ort, error) {
	// open file
	file, err := excelize.OpenFile("./unfallData/Gemeindeverzeichnis.xlsx")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// get all rows
	rows, err := file.GetRows("Onlineprodukt_Gemeinden30062026")
	if err != nil {
		return nil, err
	}

	var orte []Ort
	var landkreis string
	var gemeindeverband string

	// process data rows
	for i, row := range rows {
		// skip header row
		if i < 6 {
			continue
		}

		// ensure the row has enough columns to avoid out-of-bounds panics (names are at index 7)
		if len(row) < 7 {
			continue
		}

		// 60 = Gemeinde, 50 = Gemeindeverband, 40 = Landkreis
		satzart := strings.TrimSpace(row[0])

		// remember last values for Landkreis and Gemeindeverband
		if satzart == "40" {
			landkreis = strings.TrimSpace(row[7])
			continue
		}
		if satzart == "50" {
			gemeindeverband = strings.TrimSpace(row[7])
			continue
		}

		if satzart != "60" {
			continue
		}

		// ensure the row has enough columns to avoid out-of-bounds panics (verstädterung is at index 19)
		if len(row) < 20 {
			continue
		}

		// swap decimal commas (,) with dots (.)
		flächeStr := strings.ReplaceAll(normalizeValue(row[8]), ",", ".")
		longStr := strings.ReplaceAll(normalizeValue(row[14]), ",", ".")
		latStr := strings.ReplaceAll(normalizeValue(row[15]), ",", ".")

		// parse data
		fläche, _ := strconv.ParseFloat(flächeStr, 64)
		bevölkerung, _ := strconv.Atoi(normalizeValue(row[9]))
		männlich, _ := strconv.Atoi(normalizeValue(row[10]))
		weiblich, _ := strconv.Atoi(normalizeValue(row[11]))
		longitude, _ := strconv.ParseFloat(longStr, 64)
		latitude, _ := strconv.ParseFloat(latStr, 64)

		// map columns to the struct
		ort := Ort{
			Bundesland:         parseBundesland(row[2]),
			Regierungsbezirk:   strings.TrimSpace(row[3]),
			Kreis:              strings.TrimSpace(row[4]),
			Gemeinde:           strings.TrimSpace(row[6]),
			Name:               strings.TrimSpace(row[7]),
			Gemeindeverband:    gemeindeverband,
			Landkreis:          landkreis,
			Postleitzahl:       strings.TrimSpace(row[13]),
			Fläche:             fläche,
			Bevölkerung:        bevölkerung,
			Männlich:           männlich,
			Weiblich:           weiblich,
			Reisegebiet:        strings.TrimSpace(row[17]),
			Verstädterungsgrad: strings.TrimSpace(row[19]),
			Latitude:           latitude,
			Longitude:          longitude,
		}

		orte = append(orte, ort)
	}

	return orte, nil

}
