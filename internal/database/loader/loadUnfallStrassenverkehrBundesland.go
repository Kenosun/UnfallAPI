package loader

import (
	"database/sql"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/Kenosun/UnfallAPI/internal/parser"
)

func LoadUnfallStrassenverkehrBundesland(db *sql.DB) error {
	// load yearly data
	yearlyRecords, err := parser.ParseUnfallStrassenverkehrBundeslandYearly()
	if err != nil {
		return err
	}

	// start database transaction
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback() // safe to call, does nothing if committed

	// prepare statement
	statement, err := transaction.Prepare(`
		INSERT INTO unfall_strassenverkehr_bundesland (bundesland, strassenklasse, ortslage, jahr, anzahl)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(bundesland, strassenklasse, ortslage, jahr)
		DO UPDATE SET anzahl = excluded.anzahl
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	// helper function to insert records
	insertRecords := func(records []data.UnfallStrassenverkehrBundesland) error {
		for _, r := range records {
			_, err := statement.Exec(r.Bundesland, r.Strassenklasse, r.Ortslage, r.Jahr, r.Anzahl)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// insert records
	if err := insertRecords(yearlyRecords); err != nil {
		return err
	}

	return transaction.Commit()
}
