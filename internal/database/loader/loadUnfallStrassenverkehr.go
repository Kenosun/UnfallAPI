package loader

import (
	"database/sql"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/Kenosun/UnfallAPI/internal/parser"
)

func LoadUnfallStrassenverkehr(db *sql.DB) error {
	// load yearly data
	yearlyRecords, err := parser.ParseUnfallStrassenverkehrYearly()
	if err != nil {
		return err
	}

	// load monthly data
	monthlyRecords, err := parser.ParseUnfallStrassenverkehrMonthly()
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
		INSERT INTO unfall_strassenverkehr (strassenklasse, ortslage, kategorie, jahr, monat, anzahl)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(strassenklasse, ortslage, kategorie, jahr, monat)
		DO UPDATE SET anzahl = excluded.anzahl
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	// helper function to insert records
	insertRecords := func(records []data.UnfallStrassenverkehr) error {
		for _, r := range records {
			_, err := statement.Exec(r.Strassenklasse, r.Ortslage, r.Kategorie, r.Jahr, r.Monat, r.Anzahl)
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

	if err := insertRecords(monthlyRecords); err != nil {
		return err
	}

	return transaction.Commit()
}
