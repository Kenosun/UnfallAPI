package loader

import (
	"database/sql"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/Kenosun/UnfallAPI/internal/parser"
)

func LoadUnfallVerunglueckteBundesland(db *sql.DB) error {
	// load yearly data
	yearlyRecords, err := parser.ParseUnfallVerunglueckteBundeslandYearly()
	if err != nil {
		return err
	}

	// load monthly data
	monthlyRecords, err := parser.ParseUnfallVerunglueckteBundeslandMonthly()
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
		INSERT INTO unfall_verunglueckte_bundesland (bundesland, ortslage, schweregrad, jahr, monat, anzahl)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(bundesland, ortslage, schweregrad, jahr, monat)
		DO UPDATE SET anzahl = excluded.anzahl
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	// helper function to insert records
	insertRecords := func(records []data.UnfallVerunglueckteBundesland) error {
		for _, r := range records {
			_, err := statement.Exec(r.Bundesland, r.Ortslage, r.Schweregrad, r.Jahr, r.Monat, r.Anzahl)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// insert yearly records
	if err := insertRecords(yearlyRecords); err != nil {
		return err
	}

	// insert monthly records
	if err := insertRecords(monthlyRecords); err != nil {
		return err
	}

	return transaction.Commit()
}
