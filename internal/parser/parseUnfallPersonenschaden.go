package parser

import (
	"io"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/parser/helper"
)

type UnfallPersonenschaden struct {
	Unfalltyp   string
	Ortslage    string
	Schweregrad string
	Kategorie   string
	Jahr        int
	Monat       int // 1-12 for months, 0 for full year data
	Anzahl      int
}

func ParseUnfallPersonenschadenYearly() ([]UnfallPersonenschaden, error) {
	file, reader, err := helper.OpenCSV("./unfallData/csv/46241-0005_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallPersonenschaden
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
		if len(record) < 6 {
			continue
		}

		// identify header row
		if !headerFound && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" && record[4] == "" {
			// year columns start at index 5
			for i := 5; i < len(record); i++ {
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
			if record[0] == "" || record[1] == "" || record[2] == "" || record[3] == "" {
				continue
			}

			unfalltyp := strings.TrimSpace(record[0])
			ortslage := strings.TrimSpace(record[1])
			schweregrad := strings.TrimSpace(record[2])
			kategorie := strings.TrimSpace(record[3])

			// flatten the column values back to individual records per year
			for i, year := range years {
				colIdx := i + 5
				if colIdx >= len(record) {
					break
				}

				count, valid := helper.ParseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallPersonenschaden{
					Unfalltyp:   unfalltyp,
					Ortslage:    ortslage,
					Schweregrad: schweregrad,
					Kategorie:   kategorie,
					Jahr:        year,
					Monat:       0,
					Anzahl:      count,
				})
			}
		}
	}

	return records, nil
}

// ParseUnfallPersonenschadenMonthly parses the monthly multi-tiered header dataset (46241-0006_de.csv)
func ParseUnfallPersonenschadenMonthly() ([]UnfallPersonenschaden, error) {
	file, reader, err := helper.OpenCSV("./unfallData/csv/46241-0006_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallPersonenschaden
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

		if len(record) < 6 {
			continue
		}

		// identify year row
		if !headerFound && yearRow == nil && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" && record[4] == "" {
			// check if row actually contains years/numbers
			isYearRow := false
			for i := 5; i < len(record); i++ {
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
		if !headerFound && yearRow != nil && record[0] == "" && record[1] == "" && record[2] == "" && record[3] == "" && record[4] == "" {
			var lastValidYear int = -1

			for i := 5; i < len(record); i++ {
				yearStr := strings.TrimSpace(yearRow[i])
				monthStr := strings.ToLower(strings.TrimSpace(record[i]))

				// if a new year is explicitly listed, update lastValidYear
				if year, yErr := strconv.Atoi(yearStr); yErr == nil {
					lastValidYear = year
				}

				month := helper.ParseMonthToInt(monthStr)

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
			if record[0] == "" || record[1] == "" || record[2] == "" || record[3] == "" || strings.HasPrefix(record[0], "Tabelle") {
				continue
			}

			unfalltyp := strings.TrimSpace(record[0])
			ortslage := strings.TrimSpace(record[1])
			schweregrad := strings.TrimSpace(record[2])
			kategorie := strings.TrimSpace(record[3])

			// iterate through columns
			for i, colInfo := range columns {
				if colInfo.Year == -1 {
					continue // skip invalid/empty header columns
				}

				colIdx := i + 5
				if colIdx >= len(record) {
					break
				}

				count, valid := helper.ParseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallPersonenschaden{
					Unfalltyp:   unfalltyp,
					Ortslage:    ortslage,
					Schweregrad: schweregrad,
					Kategorie:   kategorie,
					Jahr:        colInfo.Year,
					Monat:       colInfo.Month,
					Anzahl:      count,
				})
			}
		}
	}

	return records, nil
}
