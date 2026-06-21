package handlers

import (
	"io"
	"strconv"
	"strings"
)

type UnfallStraßenverkehr struct {
	Straßenklasse string
	Ortslage      string
	Kategorie     string
	Jahr          int
	Monat         int // 1-12 for months, 0 for full year data
	Anzahl        int
}

func ParseUnfallStraßenverkehrYearly() ([]UnfallStraßenverkehr, error) {
	file, reader, err := openCSV("./unfallData/csv/46241-0003_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallStraßenverkehr
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
		if len(record) < 5 {
			continue
		}

		// identify year row
		if !headerFound && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" {
			// year columns start at index 4
			for i := 4; i < len(record); i++ {
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
			// skip footer metadata rows or incomplete blocks
			if record[0] == "" || record[1] == "" || record[2] == "" {
				continue
			}

			strassenklasse := strings.TrimSpace(record[0])
			ortslage := strings.TrimSpace(record[1])
			kategorie := strings.TrimSpace(record[2])

			// flatten the column values back to individual records per year
			for i, year := range years {
				colIdx := i + 4
				if colIdx >= len(record) {
					break
				}

				count, valid := parseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallStraßenverkehr{
					Straßenklasse: strassenklasse,
					Ortslage:      ortslage,
					Kategorie:     kategorie,
					Jahr:          year,
					Monat:         0,
					Anzahl:        count,
				})
			}
		}
	}

	return records, nil
}

func ParseUnfallStraßenverkehrMonthly() ([]UnfallStraßenverkehr, error) {
	file, reader, err := openCSV("./unfallData/csv/46241-0004_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallStraßenverkehr
	var columns []HeaderYearMonth
	var yearRow []string
	headerFound := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 5 {
			continue
		}

		// identify year row
		if !headerFound && yearRow == nil && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" {
			// check if row actually contains numbers/years
			isYearRow := false
			for i := 4; i < len(record); i++ {
				if _, err := strconv.Atoi(strings.TrimSpace(record[i])); err == nil {
					isYearRow = true
					break
				}
			}
			if isYearRow {
				yearRow = record
				continue
			}
		}

		// identify month row
		if !headerFound && yearRow != nil && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" {
			var lastValidYear int = -1

			for i := 4; i < len(record); i++ {
				yearStr := strings.TrimSpace(yearRow[i])
				monthStr := strings.ToLower(strings.TrimSpace(record[i]))

				// if a new year is explicitly listed, update lastValidYear
				if year, yErr := strconv.Atoi(yearStr); yErr == nil {
					lastValidYear = year
				}

				month := parseMonthToInt(monthStr)

				// only add valid columns where both year and month parse correctly
				if lastValidYear != -1 && month > 0 {
					columns = append(columns, HeaderYearMonth{Year: lastValidYear, Month: month})
				} else {
					// add a placeholder if columns don't align
					columns = append(columns, HeaderYearMonth{Year: -1, Month: -1})
				}
			}
			headerFound = true
			continue
		}

		// process data rows
		if headerFound {
			// skip footer metadata rows or table descriptors
			if record[0] == "" || record[1] == "" || record[2] == "" || strings.HasPrefix(record[0], "Tabelle") {
				continue
			}

			strassenklasse := strings.TrimSpace(record[0])
			ortslage := strings.TrimSpace(record[1])
			kategorie := strings.TrimSpace(record[2])

			// iterate through columns
			for i, colInfo := range columns {
				if colInfo.Year == -1 {
					continue // skip invalid/empty header columns
				}

				colIdx := i + 4
				if colIdx >= len(record) {
					break
				}

				count, valid := parseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallStraßenverkehr{
					Straßenklasse: strassenklasse,
					Ortslage:      ortslage,
					Kategorie:     kategorie,
					Jahr:          colInfo.Year,
					Monat:         colInfo.Month,
					Anzahl:        count,
				})
			}
		}
	}

	return records, nil
}
