package loader

import (
	"database/sql"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/Kenosun/UnfallAPI/internal/parser"
)

func LoadUnfallVerunglueckte(db *sql.DB) error {
	// load yearly data
	yearlyRecords, err := parser.ParseUnfallVerunglueckteYearly()
	if err != nil {
		return err
	}

	// load monthly data
	monthlyRecords, err := parser.ParseUnfallVerunglueckteMonthly()
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
		INSERT INTO unfall_verunglueckte (verkehrsart, ortslage, schweregrad, geschlecht, altersgruppe, jahr, monat, anzahl)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(verkehrsart, ortslage, schweregrad, geschlecht, altersgruppe, jahr, monat)
		DO UPDATE SET anzahl = excluded.anzahl
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	// helper function to insert records
	insertRecords := func(records []data.UnfallVerunglueckte) error {
		for _, r := range records {
			_, err := statement.Exec(r.Verkehrsart, r.Ortslage, r.Schweregrad, r.Geschlecht, r.Altersgruppe, r.Jahr, r.Monat, r.Anzahl)
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
