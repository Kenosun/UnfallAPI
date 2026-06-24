package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// initialize SQLite database and create schema tables
func InitializeDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	// enable Write-Ahead Logging mode for better performance
	_, err = db.Exec(`PRAGMA journal_mode=WAL;`)
	if err != nil {
		return nil, err
	}

	// create additional index to speed up queries
	schema := `
		CREATE TABLE IF NOT EXISTS unfall_statistik (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			unfallkategorie TEXT NOT NULL,
			ortslage TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(unfallkategorie, ortslage, jahr, monat)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_statistik_zeit ON unfall_statistik(jahr, monat);
		
		CREATE TABLE IF NOT EXISTS unfall_strassenverkehr (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			strassenklasse TEXT NOT NULL,
			ortslage TEXT NOT NULL,
			kategorie TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(strassenklasse, ortslage, kategorie, jahr, monat)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_strassenverkehr_zeit ON unfall_strassenverkehr(jahr, monat);

		CREATE TABLE IF NOT EXISTS unfall_personenschaden (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			unfalltyp TEXT NOT NULL,
			ortslage TEXT NOT NULL,
			schweregrad TEXT NOT NULL,
			kategorie TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(unfalltyp, ortslage, schweregrad, kategorie, jahr, monat)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_personenschaden_zeit ON unfall_personenschaden(jahr, monat);

		CREATE TABLE IF NOT EXISTS unfall_verunglueckte (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			verkehrsart TEXT NOT NULL,
			ortslage TEXT NOT NULL,
			schweregrad TEXT NOT NULL,
			geschlecht TEXT NOT NULL,
			altersgruppe TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(verkehrsart, ortslage, schweregrad, geschlecht, altersgruppe, jahr, monat)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_verunglueckte_zeit ON unfall_verunglueckte(jahr, monat);

		CREATE TABLE IF NOT EXISTS unfall_fehlverhalten (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			verkehrsart TEXT NOT NULL,
			fehlverhalten TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(verkehrsart, fehlverhalten, jahr, monat)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_fehlverhalten_zeit ON unfall_fehlverhalten(jahr, monat);
	
		CREATE TABLE IF NOT EXISTS unfall_beteiligung (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			verkehrsart TEXT NOT NULL,
			unfallkategorie TEXT NOT NULL,
			ortslage TEXT NOT NULL,
			geschlecht TEXT NOT NULL,
			altersgruppe TEXT NOT NULL,
			beteiligungsart TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(verkehrsart, unfallkategorie, ortslage, geschlecht, altersgruppe, beteiligungsart, jahr, monat)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_beteiligung_zeit ON unfall_beteiligung(jahr, monat);

		CREATE TABLE IF NOT EXISTS unfall_statistik_bundesland (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bundesland TEXT NOT NULL,
			unfallkategorie TEXT NOT NULL,
			ortslage TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(bundesland, unfallkategorie, ortslage, jahr, monat)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_statistik_bundesland_zeit ON unfall_statistik_bundesland(bundesland, jahr, monat);

		CREATE TABLE IF NOT EXISTS unfall_strassenverkehr_bundesland (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bundesland TEXT NOT NULL,
			strassenklasse TEXT NOT NULL,
			ortslage TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(bundesland, strassenklasse, ortslage, jahr)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_strassenverkehr_bundesland_jahr ON unfall_strassenverkehr_bundesland(bundesland, jahr);

		CREATE TABLE IF NOT EXISTS unfall_verunglueckte_bundesland (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bundesland TEXT NOT NULL,
			ortslage TEXT NOT NULL,
			schweregrad TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			anzahl INTEGER NOT NULL,
			UNIQUE(bundesland, ortslage, schweregrad, jahr, monat)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_verunglueckte_bundesland_zeit ON unfall_verunglueckte_bundesland(bundesland, jahr, monat);

		CREATE TABLE IF NOT EXISTS unfall (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bundesland TEXT NOT NULL,
			regierungsbezirk TEXT NOT NULL,
			kreis TEXT NOT NULL,
			gemeinde TEXT NOT NULL,
			jahr INTEGER NOT NULL,
			monat INTEGER NOT NULL,
			uhrzeit TEXT NOT NULL,
			wochentag TEXT NOT NULL,
			schweregrad TEXT NOT NULL,
			unfallart TEXT NOT NULL,
			unfalltyp TEXT NOT NULL,
			lichtverhaeltnis TEXT NOT NULL,
			mit_fahrrad INTEGER NOT NULL CHECK (mit_fahrrad IN (0, 1)),
			mit_pkw INTEGER NOT NULL CHECK (mit_pkw IN (0, 1)),
			mit_fussgaenger INTEGER NOT NULL CHECK (mit_fussgaenger IN (0, 1)),
			mit_kraftrad INTEGER NOT NULL CHECK (mit_kraftrad IN (0, 1)),
			mit_gueterkraftfahrzeug INTEGER NOT NULL CHECK (mit_gueterkraftfahrzeug IN (0, 1)),
			mit_sonstigen_verkehrsmittel INTEGER NOT NULL CHECK (mit_sonstigen_verkehrsmittel IN (0, 1)),
			strassenzustand TEXT NOT NULL,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			UNIQUE(jahr, monat, uhrzeit, latitude, longitude)
		);
		CREATE INDEX IF NOT EXISTS idx_unfall_ort_zeit ON unfall(bundesland, gemeinde, jahr, monat);
		CREATE INDEX IF NOT EXISTS idx_unfall_koordinaten ON unfall(latitude, longitude);
		CREATE INDEX IF NOT EXISTS idx_unfall_fahrrad ON unfall(jahr, monat, bundesland) WHERE mit_fahrrad = 1;
		CREATE INDEX IF NOT EXISTS idx_unfall_pkw ON unfall(jahr, monat, bundesland) WHERE mit_pkw = 1;
		CREATE INDEX IF NOT EXISTS idx_unfall_fussgaenger ON unfall(jahr, monat, bundesland) WHERE mit_fussgaenger = 1;
		CREATE INDEX IF NOT EXISTS idx_unfall_kraftrad ON unfall(jahr, monat, bundesland) WHERE mit_kraftrad = 1;
		CREATE INDEX IF NOT EXISTS idx_unfall_lkw ON unfall(jahr, monat, bundesland) WHERE mit_gueterkraftfahrzeug = 1;

		CREATE TABLE IF NOT EXISTS ort (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bundesland TEXT NOT NULL,
			regierungsbezirk TEXT NOT NULL,
			kreis TEXT NOT NULL,
			gemeinde TEXT NOT NULL,
			name TEXT NOT NULL,
			gemeindeverband TEXT NOT NULL,
			landkreis TEXT NOT NULL,
			postleitzahl TEXT NOT NULL,
			flaeche REAL NOT NULL,
			bevoelkerung INTEGER NOT NULL,
			maennlich INTEGER NOT NULL,
			weiblich INTEGER NOT NULL,
			reisegebiet TEXT NOT NULL,
			verstaedterungsgrad TEXT NOT NULL,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			UNIQUE(bundesland, regierungsbezirk, kreis, gemeinde)
		);
		CREATE INDEX IF NOT EXISTS idx_ort_koordinaten ON ort(latitude, longitude);
		CREATE INDEX IF NOT EXISTS idx_ort_plz ON ort(postleitzahl);
		CREATE INDEX IF NOT EXISTS idx_ort_name ON ort(name);
	`
	_, err = db.Exec(schema)
	return db, err
}
