# UnfallAPI

Eine Go-basierte API zur Verarbeitung und Abfrage von Verkehrsunfalldaten in Deutschland. Die Anwendung lädt Rohdaten automatisiert aus öffentlichen Datenquellen herunter, bereitet diese in einer lokalen SQLite-Datenbank auf und stellt die Daten über eine REST-API sowie eine interaktive Web-Demo zur Verfügung.

## Datenquellen

Die Anwendung nutzt ausschließlich öffentliche Datenquellen:

- Unfallatlas (OpenGeodata.NRW)
- Destatis (Statistisches Bundesamt)
- GENESIS-Online

## Funktionen

- Automatischer Download der Unfalldaten
- Aufbereitung und Speicherung in einer lokalen SQLite-Datenbank
- REST-API zur Abfrage von Unfall- und Statistikdaten
- Swagger/OpenAPI-Dokumentation
- Interaktive Web-Demo

## Voraussetzungen

Vor dem Start sollten folgende Voraussetzungen erfüllt sein:

- **Go (Golang):** Installierte Go-Laufzeitumgebung ([go.dev/doc/install](https://go.dev/doc/install))
- **Speicherplatz:** Mindestens **15 GB** freier Festplattenspeicher

## Installation und Start

### 1) Ins Projektverzeichnis wechseln

```bash
cd UnfallAPI
```

### 2) Abhängigkeiten installieren

```bash
go mod download
```

### 3) Anwendung starten

```bash
go run ./cmd/api/
```

### Hinweis zum ersten Start

Beim ersten Start werden automatisch:

- die Rohdaten heruntergeladen
- die SQLite-Datenbank erstellt
- die Daten importiert und aufbereitet

Dieser Vorgang kann je nach Internetverbindung und Hardware einige Zeit in Anspruch nehmen.

### Optional: Ausführbare Datei erstellen

#### Windows

```bash
go build -o unfallAPI.exe cmd/api/main.go
```

#### Linux / macOS

```bash
go build -o unfallAPI cmd/api/main.go
```

## Verfügbare Schnittstellen

- Web-Demo: `http://localhost:8080`
- REST-API: `http://localhost:8080/api/v1/`
- Swagger: `http://localhost:8080/swagger/index.html`

## Konfiguration

Zentrale Einstellungen wie `Server-Port` und `Log-Level` befinden sich in:

```text
cmd/api/main.go
```

Dort können Anpassungen vor dem Start der Anwendung vorgenommen werden.

## Daten

Die heruntergeladenen Rohdaten werden im Verzeichnis `unfallData/` gespeichert.

Je nach Datenquelle können dabei folgende Formate vorkommen:

- CSV
- TXT
- XLSX
- Shapefiles

## API-Dokumentation

Die generierten OpenAPI-/Swagger-Dateien befinden sich im Verzeichnis

```text
docs/
```

und werden für die API-Dokumentation verwendet. Dazu gehören insbesondere die Dateien `swagger.json` und `swagger.yaml`, welche die verfügbaren Endpunkte, Parameter und Antwortformate der REST-API beschreiben.
