package parser

import (
	"io"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/parser/helper"
)

type UnfallStatistikBundesland struct {
	Bundesland string
	Kategorie  string
	Ortslage   string
	Jahr       int
	Monat      int // 1-12 for months, 0 for full year data
	Anzahl     int
}

func ParseUnfallStatistikBundeslandYearly() ([]UnfallStatistikBundesland, error) {
	file, reader, err := helper.OpenCSV("./unfallData/csv/46241-0020_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallStatistikBundesland
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
			kategorie := strings.TrimSpace(record[1])
			ortslage := strings.TrimSpace(record[2])

			// flatten the column values back to individual records per year
			for i, year := range years {
				colIdx := i + 3
				if colIdx >= len(record) {
					break
				}

				count, valid := helper.ParseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallStatistikBundesland{
					Bundesland: bundesland,
					Kategorie:  kategorie,
					Ortslage:   ortslage,
					Jahr:       year,
					Monat:      0,
					Anzahl:     count,
				})
			}
		}
	}

	return records, nil
}

func ParseUnfallStatistikBundeslandMonthly() ([]UnfallStatistikBundesland, error) {
	file, reader, err := helper.OpenCSV("./unfallData/csv/46241-0021_de.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []UnfallStatistikBundesland
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

		if len(record) < 4 {
			continue
		}

		// identify year row
		if !headerFound && yearRow == nil && record[0] == "" && record[1] == "" && record[2] == "" {
			// check if row actually contains numbers/years
			isYearRow := false
			for i := 3; i < len(record); i++ {
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
		if !headerFound && yearRow != nil && record[0] == "" && record[1] == "" && record[2] == "" {
			for i := 3; i < len(record); i++ {
				yearStr := strings.TrimSpace(yearRow[i])
				monthStr := strings.ToLower(strings.TrimSpace(record[i]))

				year, yErr := strconv.Atoi(yearStr)
				month := helper.ParseMonthToInt(monthStr)

				// only add valid columns where both year and month parse correctly
				if yErr == nil && month > 0 {
					columns = append(columns, HeaderYearMonth{Year: year, Month: month})
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

			bundesland := strings.TrimSpace(record[0])
			kategorie := strings.TrimSpace(record[1])
			ortslage := strings.TrimSpace(record[2])

			// iterate through columns
			for i, colInfo := range columns {
				if colInfo.Year == -1 {
					continue // skip invalid/empty header columns
				}

				colIdx := i + 3
				if colIdx >= len(record) {
					break
				}

				count, valid := helper.ParseCount(record[colIdx])
				if !valid {
					continue
				}

				records = append(records, UnfallStatistikBundesland{
					Bundesland: bundesland,
					Kategorie:  kategorie,
					Ortslage:   ortslage,
					Jahr:       colInfo.Year,
					Monat:      colInfo.Month,
					Anzahl:     count,
				})
			}
		}
	}

	return records, nil
}
