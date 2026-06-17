package handlers

import (
	"io"
	"strconv"
	"strings"
)

type UnfallStraßenverkehrBundesland struct {
	Bundesland    string
	Straßenklasse string
	Ortslage      string
	Jahr          int
	Monat         int // 1-12 for months, 0 for full year data
	Anzahl        int
}

func ParseUnfallStraßenverkehrBundeslandYearly() ([]UnfallStraßenverkehrBundesland, error) {
	file, reader, err := openCSV("./unfallData/csv/46241-0022_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallStraßenverkehrBundesland
	var years []int
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
		if len(record) < 4 {
			continue
		}

		// identify year row
		if !headerFound && record[0] == "" && record[1] == "" && record[2] == "" {
			for i := 3; i < len(record); i++ {
				yearStr := strings.TrimSpace(record[i])
				if yearStr == "" {
					continue
				}
				year, err := strconv.Atoi(yearStr)
				if err == nil {
					years = append(years, year)
				}
			}
			if len(years) > 0 {
				headerFound = true
			}
			continue
		}

		// process data rows
		if headerFound {
			// skip footer metadata rows or table descriptors
			if record[0] == "" || record[1] == "" || record[2] == "" || strings.HasPrefix(record[0], "Tabelle") {
				continue
			}

			bundesland := strings.TrimSpace(record[0])
			strassenklasse := strings.TrimSpace(record[1])
			ortslage := strings.TrimSpace(record[2])

			// flatten the column values back to individual records per year
			for i, year := range years {
				colIdx := i + 3
				if colIdx >= len(record) {
					break
				}

				count, valid := parseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallStraßenverkehrBundesland{
					Bundesland:    bundesland,
					Straßenklasse: strassenklasse,
					Ortslage:      ortslage,
					Jahr:          year,
					Monat:         0,
					Anzahl:        count,
				})
			}
		}
	}

	return records, nil
}
