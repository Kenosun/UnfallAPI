package loader

import (
	"database/sql"

	"github.com/Kenosun/UnfallAPI/internal/parser"
)

func LoadOrt(db *sql.DB) error {
	// load data
	ortRecords, err := parser.ParseOrt()
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
		INSERT INTO ort (
			bundesland, regierungsbezirk, kreis, gemeinde, name, 
			gemeindeverband, landkreis, postleitzahl, flaeche, bevoelkerung, 
			maennlich, weiblich, reisegebiet, verstaedterungsgrad, latitude, longitude
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(bundesland, regierungsbezirk, kreis, gemeinde)
		DO UPDATE SET
			name = excluded.name,
			gemeindeverband = excluded.gemeindeverband,
			landkreis = excluded.landkreis,
			postleitzahl = excluded.postleitzahl,
			flaeche = excluded.flaeche,
			bevoelkerung = excluded.bevoelkerung,
			maennlich = excluded.maennlich,
			weiblich = excluded.weiblich,
			reisegebiet = excluded.reisegebiet,
			verstaedterungsgrad = excluded.verstaedterungsgrad,
			latitude = excluded.latitude,
			longitude = excluded.longitude
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	// insert records
	for _, r := range ortRecords {
		_, err := statement.Exec(
			r.Bundesland,
			r.Regierungsbezirk,
			r.Kreis,
			r.Gemeinde,
			r.Name,
			r.Gemeindeverband,
			r.Landkreis,
			r.Postleitzahl,
			r.Flaeche,
			r.Bevoelkerung,
			r.Maennlich,
			r.Weiblich,
			r.Reisegebiet,
			r.Verstaedterungsgrad,
			r.Latitude,
			r.Longitude,
		)
		if err != nil {
			return err
		}
	}

	return transaction.Commit()
}
