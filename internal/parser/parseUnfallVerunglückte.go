package parser

import (
	"io"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/parser/helper"
)

type UnfallVerunglückte struct {
	Verkehrsart  string
	Ortslage     string
	Kategorie    string
	Geschlecht   string
	Altersgruppe string
	Jahr         int
	Monat        int // 1-12 for months, 0 for full year data
	Anzahl       int
}

func ParseUnfallVerunglückteYearly() ([]UnfallVerunglückte, error) {
	file, reader, err := helper.OpenCSV("./unfallData/csv/46241-0007_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallVerunglückte
	var headers []HeaderGenderAge
	var geschlechtRow []string
	headerFound := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// skip empty or incomplete metadata rows
		if len(record) < 5 {
			continue
		}

		// identify gender row
		if !headerFound && geschlechtRow == nil && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" {
			if strings.Contains(record[4], "Geschlecht") {
				geschlechtRow = record
				continue
			}
		}

		// set gender row to sub-gender values
		if geschlechtRow != nil && headers == nil && (record[4] == "männlich" || record[4] == "weiblich" || record[4] == "Insgesamt") {
			geschlechtRow = record
			continue
		}

		// identify age row
		if !headerFound && geschlechtRow != nil && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" {
			var lastValidGeschlecht string
			// numerical columns start at index 4
			for i := 4; i < len(record); i++ {
				gStr := strings.TrimSpace(geschlechtRow[i])
				aStr := strings.TrimSpace(record[i])

				if gStr != "" && gStr != "Geschlecht" {
					lastValidGeschlecht = gStr
				}

				if lastValidGeschlecht != "" && aStr != "" {
					headers = append(headers, HeaderGenderAge{
						Geschlecht:   lastValidGeschlecht,
						Altersgruppe: aStr,
					})
				} else {
					headers = append(headers, HeaderGenderAge{Geschlecht: "Unknown", Altersgruppe: "Unknown"})
				}
			}
			headerFound = true
			continue
		}

		// process data rows
		if headerFound {
			// skip footer metadata rows or table descriptors
			if record[0] == "" || record[1] == "" || record[2] == "" || record[3] == "" || strings.HasPrefix(record[0], "Tabelle") {
				continue
			}

			// year is on every row in column 0
			year, err := strconv.Atoi(strings.TrimSpace(record[0]))
			if err != nil {
				continue // skip row if year is invalid
			}

			verkehrsart := strings.TrimSpace(record[1])
			ortslage := strings.TrimSpace(record[2])
			kategorie := strings.TrimSpace(record[3])

			// flatten columns using headers slice
			for i, header := range headers {
				colIdx := i + 4
				if colIdx >= len(record) {
					break
				}

				count, valid := helper.ParseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallVerunglückte{
					Verkehrsart:  verkehrsart,
					Ortslage:     ortslage,
					Kategorie:    kategorie,
					Geschlecht:   header.Geschlecht,
					Altersgruppe: header.Altersgruppe,
					Jahr:         year,
					Monat:        0,
					Anzahl:       count,
				})
			}
		}
	}

	return records, nil
}

func ParseUnfallVerunglückteMonthly() ([]UnfallVerunglückte, error) {
	file, reader, err := helper.OpenCSV("./unfallData/csv/46241-0008_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallVerunglückte
	var headers []HeaderGenderAge
	var geschlechtRow []string
	headerFound := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// skip empty or incomplete metadata rows
		if len(record) < 6 {
			continue
		}

		// identify gender row
		if !headerFound && geschlechtRow == nil && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" && record[4] == "" {
			if strings.Contains(record[5], "Geschlecht") {
				geschlechtRow = record
				continue
			}
		}

		// set gender row to sub-gender values
		if geschlechtRow != nil && headers == nil && (record[5] == "männlich" || record[5] == "weiblich" || record[5] == "Insgesamt") {
			geschlechtRow = record
			continue
		}

		// identify age row
		if !headerFound && geschlechtRow != nil && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" && record[4] == "" {
			var lastValidGeschlecht string
			// numerical columns start at index 5
			for i := 5; i < len(record); i++ {
				gStr := strings.TrimSpace(geschlechtRow[i])
				aStr := strings.TrimSpace(record[i])

				if gStr != "" && gStr != "Geschlecht" {
					lastValidGeschlecht = gStr
				}

				if lastValidGeschlecht != "" && aStr != "" {
					headers = append(headers, HeaderGenderAge{
						Geschlecht:   lastValidGeschlecht,
						Altersgruppe: aStr,
					})
				} else {
					headers = append(headers, HeaderGenderAge{Geschlecht: "Unknown", Altersgruppe: "Unknown"})
				}
			}
			headerFound = true
			continue
		}

		// process data rows
		if headerFound {
			// skip footer metadata rows or table descriptors
			if record[0] == "" || record[1] == "" || record[2] == "" || record[3] == "" || record[4] == "" || strings.HasPrefix(record[0], "Tabelle") {
				continue
			}

			// year is on every row in column 0
			year, err := strconv.Atoi(strings.TrimSpace(record[0]))
			if err != nil {
				continue // skip row if year is invalid
			}

			monthStr := strings.ToLower(strings.TrimSpace(record[1]))
			month := helper.ParseMonthToInt(monthStr)
			if month == -1 {
				continue // skip row if month is invalid
			}

			verkehrsart := strings.TrimSpace(record[2])
			ortslage := strings.TrimSpace(record[3])
			kategorie := strings.TrimSpace(record[4])

			// flatten columns using headers slice
			for i, header := range headers {
				colIdx := i + 5
				if colIdx >= len(record) {
					break
				}

				count, valid := helper.ParseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallVerunglückte{
					Verkehrsart:  verkehrsart,
					Ortslage:     ortslage,
					Kategorie:    kategorie,
					Geschlecht:   header.Geschlecht,
					Altersgruppe: header.Altersgruppe,
					Jahr:         year,
					Monat:        month,
					Anzahl:       count,
				})
			}
		}
	}

	return records, nil
}
