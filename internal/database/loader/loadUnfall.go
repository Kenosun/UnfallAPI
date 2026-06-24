package loader

import (
	"database/sql"
	"path/filepath"

	"github.com/Kenosun/UnfallAPI/internal/parser"
)

func LoadUnfall(db *sql.DB) error {
	// get all files matching the pattern
	csvFiles, err := filepath.Glob("./unfallData/csv/Unfallort*.csv")
	if err != nil {
		return err
	}

	txtFiles, err := filepath.Glob("./unfallData/csv/Unfallort*.txt")
	if err != nil {
		return err
	}

	files := append(csvFiles, txtFiles...)

	// start a single database transaction for all files
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback() // safe to call, does nothing if committed

	// prepare statement
	statement, err := transaction.Prepare(`
		INSERT INTO unfall (
			bundesland, regierungsbezirk, kreis, gemeinde, jahr, monat, uhrzeit, wochentag, 
			schweregrad, unfallart, unfalltyp, lichtverhaeltnis, mit_fahrrad, mit_pkw, 
			mit_fussgaenger, mit_kraftrad, mit_gueterkraftfahrzeug, mit_sonstigen_verkehrsmittel, 
			strassenzustand, latitude, longitude
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(jahr, monat, uhrzeit, latitude, longitude)
		DO UPDATE SET 
			bundesland = excluded.bundesland,
			regierungsbezirk = excluded.regierungsbezirk,
			kreis = excluded.kreis,
			gemeinde = excluded.gemeinde,
			wochentag = excluded.wochentag,
			schweregrad = excluded.schweregrad,
			unfallart = excluded.unfallart,
			unfalltyp = excluded.unfalltyp,
			lichtverhaeltnis = excluded.lichtverhaeltnis,
			mit_fahrrad = excluded.mit_fahrrad,
			mit_pkw = excluded.mit_pkw,
			mit_fussgaenger = excluded.mit_fussgaenger,
			mit_kraftrad = excluded.mit_kraftrad,
			mit_gueterkraftfahrzeug = excluded.mit_gueterkraftfahrzeug,
			mit_sonstigen_verkehrsmittel = excluded.mit_sonstigen_verkehrsmittel,
			strassenzustand = excluded.strassenzustand
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	// loop through each file
	for _, filePath := range files {
		// load data
		records, err := parser.ParseUnfall(filePath)
		if err != nil {
			return err
		}

		// insert records
		for _, r := range records {
			_, err := statement.Exec(
				r.Bundesland, r.Regierungsbezirk, r.Kreis, r.Gemeinde, r.Jahr, r.Monat, r.Uhrzeit, r.Wochentag,
				r.Schweregrad, r.Unfallart, r.Unfalltyp, r.Lichtverhaeltnis, r.MitFahrrad, r.MitPKW,
				r.MitFussgaenger, r.MitKraftrad, r.MitGueterkraftfahrzeug, r.MitSonstigenVerkehrsmittel,
				r.Strassenzustand, r.Latitude, r.Longitude,
			)
			if err != nil {
				return err
			}
		}
	}

	return transaction.Commit()
}
