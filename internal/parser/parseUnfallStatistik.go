package parser

import (
	"io"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/Kenosun/UnfallAPI/internal/parser/helper"
)

func ParseUnfallStatistikYearly() ([]data.UnfallStatistik, error) {
	file, reader, err := helper.OpenCSV("./unfallData/csv/46241-0001_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []data.UnfallStatistik
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
		if len(record) < 3 {
			continue
		}

		// identify year row
		if !headerFound && record[0] == "" && record[1] == "" {
			for i := 2; i < len(record); i++ {
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
			if record[0] == "" || record[1] == "" {
				continue
			}

			kategorie := strings.TrimSpace(record[0])
			ortslage := strings.TrimSpace(record[1])

			// flatten the column values back to individual records per year
			for i, year := range years {
				colIdx := i + 2
				if colIdx >= len(record) {
					break
				}

				count, valid := helper.ParseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, data.UnfallStatistik{
					Kategorie: kategorie,
					Ortslage:  ortslage,
					Jahr:      year,
					Monat:     0,
					Anzahl:    count,
				})
			}
		}
	}

	return records, nil
}

func ParseUnfallStatistikMonthly() ([]data.UnfallStatistik, error) {
	file, reader, err := helper.OpenCSV("./unfallData/csv/46241-0002_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []data.UnfallStatistik
	var columns []data.HeaderYearMonth
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

		if len(record) < 3 {
			continue
		}

		// identify year row
		if !headerFound && yearRow == nil && record[0] == "" && record[1] == "" {
			// check if row actually contains numbers/years
			isYearRow := false
			for i := 2; i < len(record); i++ {
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
		if !headerFound && yearRow != nil && record[0] == "" && record[1] == "" {
			for i := 2; i < len(record); i++ {
				yearStr := strings.TrimSpace(yearRow[i])
				monthStr := strings.ToLower(strings.TrimSpace(record[i]))

				year, yErr := strconv.Atoi(yearStr)
				month := helper.ParseMonthToInt(monthStr)

				// only add valid columns where both year and month parse correctly
				if yErr == nil && month > 0 {
					columns = append(columns, data.HeaderYearMonth{Year: year, Month: month})
				} else {
					// add a placeholder if columns don't align
					columns = append(columns, data.HeaderYearMonth{Year: -1, Month: -1})
				}
			}
			headerFound = true
			continue
		}

		// process data rows
		if headerFound {
			// skip footer metadata rows or table descriptors
			if record[0] == "" || record[1] == "" || strings.HasPrefix(record[0], "Tabelle") {
				continue
			}

			kategorie := strings.TrimSpace(record[0])
			ortslage := strings.TrimSpace(record[1])

			// iterate through columns
			for i, colInfo := range columns {
				if colInfo.Year == -1 {
					continue // skip invalid/empty header columns
				}

				colIdx := i + 2
				if colIdx >= len(record) {
					break
				}

				count, valid := helper.ParseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, data.UnfallStatistik{
					Kategorie: kategorie,
					Ortslage:  ortslage,
					Jahr:      colInfo.Year,
					Monat:     colInfo.Month,
					Anzahl:    count,
				})
			}
		}
	}

	return records, nil
}
