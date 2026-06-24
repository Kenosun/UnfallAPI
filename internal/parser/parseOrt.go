package parser

import (
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/Kenosun/UnfallAPI/internal/parser/helper"
	"github.com/xuri/excelize/v2"
)

func normalizeString(s string) string {
	s = strings.ReplaceAll(s, " ", "")      // remove spaces
	s = strings.ReplaceAll(s, "\u00a0", "") // handle non-breaking spaces
	return s
}

func parseInt(s string) int {
	val, err := strconv.Atoi(normalizeString(s))
	if err != nil {
		return -1
	}
	return val
}

func parseFloat(s string) float64 {
	val, err := strconv.ParseFloat(normalizeString(s), 64)
	if err != nil {
		return -1.0
	}
	return val
}

func ParseOrt() ([]data.Ort, error) {
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

	var orte []data.Ort
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

		// ensure the row has enough columns to avoid out-of-bounds panics (verstaedterung is at index 19)
		if len(row) < 20 {
			continue
		}

		// swap decimal commas (,) with dots (.)
		flaecheStr := strings.ReplaceAll(normalizeString(row[8]), ",", ".")
		longitudeStr := strings.ReplaceAll(normalizeString(row[14]), ",", ".")
		latitudeStr := strings.ReplaceAll(normalizeString(row[15]), ",", ".")

		// parse data
		flaeche := parseFloat(flaecheStr)
		bevoelkerung := parseInt((row[9]))
		maennlich := parseInt((row[10]))
		weiblich := parseInt((row[11]))
		longitude := parseFloat(longitudeStr)
		latitude := parseFloat(latitudeStr)

		// map columns to the struct
		ort := data.Ort{
			Bundesland:          helper.ParseBundesland(row[2]),
			Regierungsbezirk:    strings.TrimSpace(row[3]),
			Kreis:               strings.TrimSpace(row[4]),
			Gemeinde:            strings.TrimSpace(row[6]),
			Name:                strings.TrimSpace(row[7]),
			Gemeindeverband:     gemeindeverband,
			Landkreis:           landkreis,
			Postleitzahl:        strings.TrimSpace(row[13]),
			Flaeche:             flaeche,
			Bevoelkerung:        bevoelkerung,
			Maennlich:           maennlich,
			Weiblich:            weiblich,
			Reisegebiet:         strings.TrimSpace(row[17]),
			Verstaedterungsgrad: strings.TrimSpace(row[19]),
			Latitude:            latitude,
			Longitude:           longitude,
		}

		orte = append(orte, ort)
	}

	return orte, nil

}
