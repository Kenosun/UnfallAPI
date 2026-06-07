package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const (
	destDir  = "./unfallData"
	tempFile = destDir + "/download.zip"
	csvPath  = destDir + "/csv"
	shpPath  = destDir + "/shp"
)

var urls = [...]string{
	// Straßenverkehrsunfälle OpenGeodata NRW -> Unfallatlas
	// 2024
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2024_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2024_EPSG25832_Shape.zip",
	// 2023
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2023_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2023_EPSG25832_Shape.zip",
	// 2022
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2022_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2022_EPSG25832_Shape.zip",
	// 2021
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2021_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2021_EPSG25832_Shape.zip",
	// 2020
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2020_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2020_EPSG25832_Shape.zip",
	// 2019
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2019_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2019_EPSG25832_Shape.zip",
	// 2018
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2018_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2018_EPSG25832_Shape.zip",
	// 2017
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2017_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2017_EPSG25832_Shape.zip",
	// 2016
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2016_EPSG25832_CSV.zip",
	"https://www.opengeodata.nrw.de/produkte/transport_verkehr/unfallatlas/Unfallorte2016_EPSG25832_Shape.zip",
	// Regionalatlas
	// Straßenverkehrsunfälle bezogen auf EW
	"https://www.regionalstatistik.de:443/genesis/online?operation=ergebnistabelleDownload&levelindex=1&levelid=1780507533219&option=dcsv",
	// Straßenverkehrsunfälle bezogen auf Kfz
	"https://www.regionalstatistik.de:443/genesis/online?operation=ergebnistabelleDownload&levelindex=1&levelid=1780507951673&option=dcsv",
	// Getötete bei Straßenverkehrsunfällen je 100.000 EW
	"https://www.regionalstatistik.de:443/genesis/online?operation=ergebnistabelleDownload&levelindex=1&levelid=1780508058721&option=dcsv",
	// Gemeindeverzeichnis
	// Quartalsausgabe als .xlsx
	"https://www.destatis.de/DE/Themen/Laender-Regionen/Regionales/Gemeindeverzeichnis/Administrativ/Archiv/GVAuszugQ/AuszugGV2QAktuell.xlsx?__blob=publicationFile&v=13",
}

func LoadUnfallData() error {
	// data already exists
	_, errCsv := os.Stat(csvPath)
	_, errShp := os.Stat(shpPath)
	if errCsv == nil && errShp == nil {
		return nil
	}

	if err := os.MkdirAll(csvPath, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(shpPath, 0755); err != nil {
		return err
	}

	// progress bar
	bar := progressbar.Default(int64(len(urls)), "Downloading UnfallData...")

	// download and extract loop
	for _, url := range urls {
		if err := download(url, tempFile); err != nil {
			return fmt.Errorf("error downloading %s: %w", url, err)
		}

		// extract zip files
		if strings.Contains(strings.ToLower(url), ".zip") {
			if err := extract(tempFile, destDir); err != nil {
				_ = os.Remove(tempFile)
				return fmt.Errorf("error extracting %s: %w", url, err)
			}
		} else {
			// move raw files directly instead of unzipping (hardcoded for now)
			var finalPath string
			if strings.Contains(url, "xlsx") {
				finalPath = filepath.Join(destDir, "Gemeindeverzeichnis.xlsx")
			}
			_ = os.Rename(tempFile, finalPath)
		}

		_ = os.Remove(tempFile)
		_ = bar.Add(1)
	}

	// copy genesisData to csvPath
	if err := extract("./genesisData.zip", csvPath); err != nil {
		return fmt.Errorf("error copying genesisData: %w", err)
	}

	return nil
}

func download(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	return out.Sync()
}

func extract(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		// prevent directory traversal
		cleanedPath := filepath.Clean(f.Name)
		if strings.HasPrefix(cleanedPath, "..") || strings.HasPrefix(cleanedPath, "/") {
			continue
		}

		// nested zip -> extract it temporarily and recurse into it
		if strings.ToLower(filepath.Ext(f.Name)) == ".zip" {
			nestedZipPath := filepath.Join(os.TempDir(), filepath.Base(f.Name))

			if err := extractFile(f, nestedZipPath); err != nil {
				return err
			}

			if err := extract(nestedZipPath, destDir); err != nil {
				_ = os.Remove(nestedZipPath)
				return err
			}

			_ = os.Remove(nestedZipPath)
			continue
		}

		// normal file (CSV, SHP, DBF, etc.) -> route directly to its final destination
		ext := strings.ToLower(filepath.Ext(f.Name))
		var targetDir string
		if ext == ".csv" || ext == ".txt" || ext == ".ini" {
			targetDir = csvPath
		} else {
			targetDir = shpPath
		}

		targetPath := filepath.Join(targetDir, filepath.Base(f.Name))
		if err := extractFile(f, targetPath); err != nil {
			return err
		}
	}
	return nil
}

func extractFile(f *zip.File, targetPath string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
