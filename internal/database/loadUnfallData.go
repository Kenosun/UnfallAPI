package database

import (
	"database/sql"

	"github.com/Kenosun/UnfallAPI/internal/database/loader"
	"github.com/charmbracelet/log"
	"github.com/schollz/progressbar/v3"
)

func LoadUnfallData(db *sql.DB) error {
	// define loader functions and their names for error logging
	tasks := []struct {
		name string
		run  func(*sql.DB) error
	}{
		{"UnfallStatistik", loader.LoadUnfallStatistik},
		{"UnfallStrassenverkehr", loader.LoadUnfallStrassenverkehr},
		{"UnfallPersonenschaden", loader.LoadUnfallPersonenschaden},
		{"UnfallVerunglueckte", loader.LoadUnfallVerunglueckte},
		{"UnfallFehlverhalten", loader.LoadUnfallFehlverhalten},
		{"UnfallBeteiligung", loader.LoadUnfallBeteiligung},
		{"UnfallStatistikBundesland", loader.LoadUnfallStatistikBundesland},
		{"UnfallStrassenverkehrBundesland", loader.LoadUnfallStrassenverkehrBundesland},
		{"UnfallVerunglueckteBundesland", loader.LoadUnfallVerunglueckteBundesland},
		{"Unfall (Unfallatlas)", loader.LoadUnfall},
		{"Ort (Gemeindeverzeichnis)", loader.LoadOrt},
	}

	// initialize progress bar
	bar := progressbar.Default(int64(len(tasks)), "Loading UnfallData into Database...")

	// iterate through tasks and increment the bar
	for _, task := range tasks {
		bar.Describe("Loading " + task.name)

		if err := task.run(db); err != nil {
			log.Error("Failed to load %s: %v", task.name, err)
		}

		_ = bar.Add(1)
	}

	return nil
}
