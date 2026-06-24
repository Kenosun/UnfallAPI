package loader

import (
	"database/sql"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/Kenosun/UnfallAPI/internal/parser"
)

func LoadUnfallBeteiligung(db *sql.DB) error {
	// load yearly data
	yearlyRecords, err := parser.ParseUnfallBeteiligungYearly()
	if err != nil {
		return err
	}

	// load monthly data
	monthlyRecords, err := parser.ParseUnfallBeteiligungMonthly()
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
		INSERT INTO unfall_beteiligung (verkehrsart, unfallkategorie, ortslage, geschlecht, altersgruppe, beteiligungsart, jahr, monat, anzahl)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(verkehrsart, unfallkategorie, ortslage, geschlecht, altersgruppe, beteiligungsart, jahr, monat)
		DO UPDATE SET anzahl = excluded.anzahl
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	// helper function to insert records
	insertRecords := func(records []data.UnfallBeteiligung) error {
		for _, r := range records {
			_, err := statement.Exec(r.Verkehrsart, r.Unfallkategorie, r.Ortslage, r.Geschlecht, r.Altersgruppe, r.Beteiligungsart, r.Jahr, r.Monat, r.Anzahl)
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
