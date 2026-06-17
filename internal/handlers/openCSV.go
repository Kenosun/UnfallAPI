package handlers

import (
	"encoding/csv"
	"os"
)

func openCSV(path string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'          // set delimiter to semicolon
	reader.FieldsPerRecord = -1 // turn off strict fields per record checking
	return file, reader, nil
}
