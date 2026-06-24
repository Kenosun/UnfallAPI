package parser

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/Kenosun/UnfallAPI/internal/parser/helper"
)

func ParseUnfall(filePath string) ([]data.Unfall, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	// identify header row
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// map uppercase header names to indices
	headerMap := make(map[string]int)
	for idx, name := range header {
		cleanName := strings.TrimSpace(name)
		cleanName = strings.TrimPrefix(cleanName, "\ufeff")
		headerMap[strings.ToUpper(cleanName)] = idx
	}

	var records []data.Unfall

	// process rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// skip empty rows
		if len(record) == 0 {
			continue
		}

		// parse string field by its uppercase header name
		parseString := func(fieldName string) string {
			if idx, ok := headerMap[fieldName]; ok && idx < len(record) {
				// special case for Uhrzeit
				if fieldName == "USTUNDE" {
					return strings.TrimSpace(record[idx]) + ":00 Uhr"
				}
				return strings.TrimSpace(record[idx])
			}
			return ""
		}

		// parse integer field (defaults to -1 on failure)
		parseInt := func(fieldName string) int {
			if idx, ok := headerMap[fieldName]; ok && idx < len(record) {
				val, err := strconv.Atoi(strings.TrimSpace(record[idx]))
				if err != nil {
					return -1
				}
				return val
			}
			return -1
		}

		// convert binary indicators into booleans
		parseBool := func(fieldName string) bool {
			if idx, ok := headerMap[fieldName]; ok && idx < len(record) {
				return strings.TrimSpace(record[idx]) == "1"
			}
			return false
		}

		// parse floats by swapping decimal commas (,) with dots (.)
		parseFloat := func(fieldName string) float64 {
			if idx, ok := headerMap[fieldName]; ok && idx < len(record) {
				valStr := strings.ReplaceAll(strings.TrimSpace(record[idx]), ",", ".")
				val, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					return -1.0
				}
				return val
			}
			return -1.0
		}

		// parse integer field with multiple header variations
		parseIntWithHeaderVariation := func(fields ...string) int {
			for _, field := range fields {
				if _, ok := headerMap[field]; ok {
					return parseInt(field)
				}
			}
			return -1
		}

		// parse bool field with multiple header variations
		parseBoolWithHeaderVariation := func(fields ...string) bool {
			for _, field := range fields {
				if _, ok := headerMap[field]; ok {
					return parseBool(field)
				}
			}
			return false
		}

		item := data.Unfall{
			Bundesland:                 helper.ParseBundesland(parseString("ULAND")),
			Regierungsbezirk:           parseString("UREGBEZ"),
			Kreis:                      parseString("UKREIS"),
			Gemeinde:                   parseString("UGEMEINDE"),
			Jahr:                       parseInt("UJAHR"),
			Monat:                      parseInt("UMONAT"),
			Uhrzeit:                    parseString("USTUNDE"),
			Wochentag:                  helper.ParseWochentag(parseInt("UWOCHENTAG")),
			Schweregrad:                helper.ParseSchweregrad(parseInt("UKATEGORIE")),
			Unfallart:                  helper.ParseUnfallart(parseInt("UART")),
			Unfalltyp:                  helper.ParseUnfalltyp(parseInt("UTYP1")),
			Lichtverhaeltnis:           helper.ParseLichtverhaeltnis(parseIntWithHeaderVariation("ULICHTVERH", "LICHT")),
			MitFahrrad:                 parseBool("ISTRAD"),
			MitPKW:                     parseBool("ISTPKW"),
			MitFussgaenger:             parseBool("ISTFUSS"),
			MitKraftrad:                parseBool("ISTKRAD"),
			MitGueterkraftfahrzeug:     parseBool("ISTGKFZ"),
			MitSonstigenVerkehrsmittel: parseBoolWithHeaderVariation("ISTSONSTIGE", "ISTSONSTIG"),
			Strassenzustand:            helper.ParseStrassenzustand(parseIntWithHeaderVariation("STRZUSTAND", "ISTSTRASSENZUSTAND", "ISTSTRASSE")),
			Latitude:                   parseFloat("YGCSWGS84"),
			Longitude:                  parseFloat("XGCSWGS84"),
		}

		records = append(records, item)
	}

	return records, nil
}
