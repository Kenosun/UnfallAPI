package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/gin-gonic/gin"
)

// GetUnfall godoc
//
// @Summary      Detalierte Verkehrsunfälle mit Personenschaden im Unfallatlas abrufen
// @Description  Detaillierte, punktgenaue Daten zu Straßenverkehrsunfällen inklusive Geokoordinaten und Unfallumständen abrufen
// @Tags         Unfallatlas
// @Accept       json
// @Produce      json
// @Param        bundesland       				query    string  false  "Filter nach Bundesland" enums(Baden-Württemberg, Bayern, Berlin, Brandenburg, Bremen, Hamburg, Hessen, Mecklenburg-Vorpommern, Niedersachsen, Nordrhein-Westfalen, Rheinland-Pfalz, Saarland, Sachsen, Sachsen-Anhalt, Schleswig-Holstein, Thüringen)
// @Param        regierungsbezirk 				query    string  false  "Filter nach Regierungsbezirk (Amtlicher Gemeindeschlüssel)"
// @Param        kreis            				query    string  false  "Filter nach Kreis (Amtlicher Gemeindeschlüssel)"
// @Param        gemeinde         				query    string  false  "Filter nach Gemeinde (Amtlicher Gemeindeschlüssel)"
// @Param        jahr             				query    int     false  "Filter nach Jahr"
// @Param        monat            				query    int     false  "Filter nach Monat (1-12)" minimum(1) maximum(12)
// @Param        uhrzeit          				query    string  false  "Filter nach Uhrzeit (xx:00 Uhr)"
// @Param        wochentag        				query    string  false  "Filter nach Wochentag" enums(Montag, Dienstag, Mittwoch, Donnerstag, Freitag, Samstag, Sonntag)
// @Param        schweregrad      				query    string  false  "Filter nach Schweregrad" enums(Unfall mit Getöteten, Unfall mit Schwerverletzten, Unfall mit Leichtverletzten)
// @Param        unfallart        				query    string  false  "Filter nach Unfallart" enums(Zusammenstoß mit anfahrendem/anhaltendem/ruhendem Fahrzeug, Zusammenstoß mit vorausfahrendem/wartendem Fahrzeug, Zusammenstoß mit seitlich in gleicher Richtung fahrendem Fahrzeug, Zusammenstoß mit entgegenkommendem Fahrzeug, Zusammenstoß mit einbiegendem/kreuzendem Fahrzeug, Zusammenstoß zwischen Fahrzeug und Fußgänger, Aufprall auf Fahrbahnhindernis, Abkommen von Fahrbahn nach rechts, Abkommen von Fahrbahn nach links, Unfall anderer Art)
// @Param        unfalltyp        				query    string  false  "Filter nach Unfalltyp" enums(Fahrunfall, Abbiegeunfall, Einbiegen/Kreuzen-Unfall, Überschreiten-Unfall, Unfall durch ruhenden Verkehr, Unfall im Längsverkehr, sonstiger Unfall)
// @Param        lichtverhaeltnis 				query    string  false  "Filter nach Lichtverhätltnis" enums(Tageslicht, Dämmerung, Dunkelheit)
// @Param        mit_fahrrad      				query    bool    false  "Beteiligung von Fahrrädern (true/false)"
// @Param        mit_pkw          				query    bool    false  "Beteiligung von PKWs (true/false)"
// @Param        mit_fussgaenger  				query    bool    false  "Beteiligung von Fußgängern (true/false)"
// @Param        mit_kraftrad     				query    bool    false  "Beteiligung von Kraftfahrzeugen (true/false)"
// @Param        mit_gueterkraftfahrzeug 		query    bool    false  "Beteiligung von Güterkraftfahrzeugen (true/false)"
// @Param        mit_sonstigen_verkehrsmittel 	query    bool    false  "Beteiligung von sonstigen Verkehrsmitteln (true/false)"
// @Param        strassenzustand  				query    string  false  "Filter nach Straßenzustand" enums(trocken, nass/feucht/schlüpfrig, winterglatt)
// @Param        min_lat          				query    float64 false  "Minimum Latitude"
// @Param        max_lat          				query    float64 false  "Maximum Latitude"
// @Param        min_lon          				query    float64 false  "Minimum Longitude"
// @Param        max_lon          				query    float64 false  "Maximum Longitude"
// @Success      200              				{array}  data.Unfall
// @Failure      400         					{object} HTTPError "Bad Request - Invalid parameter type or range"
// @Failure      500         					{object} HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfall [get]
func (h *AccidentHandler) GetUnfall(c *gin.Context) {
	baseQuery := `
		SELECT 
			bundesland, regierungsbezirk, kreis, gemeinde, jahr, monat, uhrzeit, wochentag, 
			schweregrad, unfallart, unfalltyp, lichtverhaeltnis, mit_fahrrad, mit_pkw, 
			mit_fussgaenger, mit_kraftrad, mit_gueterkraftfahrzeug, mit_sonstigen_verkehrsmittel,
			strassenzustand, latitude, longitude 
		FROM unfall`

	// dynamic WHERE clauses
	var whereClauses []string
	var queryArgs []any

	// string parameters
	stringParams := map[string]string{
		"bundesland":       c.Query("bundesland"),
		"regierungsbezirk": c.Query("regierungsbezirk"),
		"kreis":            c.Query("kreis"),
		"gemeinde":         c.Query("gemeinde"),
		"uhrzeit":          c.Query("uhrzeit"),
		"wochentag":        c.Query("wochentag"),
		"schweregrad":      c.Query("schweregrad"),
		"unfallart":        c.Query("unfallart"),
		"unfalltyp":        c.Query("unfalltyp"),
		"lichtverhaeltnis": c.Query("lichtverhaeltnis"),
		"strassenzustand":  c.Query("strassenzustand"),
	}

	for column, value := range stringParams {
		if value != "" {
			whereClauses = append(whereClauses, column+" = ?")
			queryArgs = append(queryArgs, value)
		}
	}

	// boolean parameters
	boolParams := []string{
		"mit_fahrrad", "mit_pkw", "mit_fussgaenger", "mit_kraftrad",
		"mit_gueterkraftfahrzeug", "mit_sonstigen_verkehrsmittel",
	}

	for _, param := range boolParams {
		if valStr := c.Query(param); valStr != "" {
			val, err := strconv.ParseBool(valStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, HTTPError{Error: "invalid " + param + " format (must be true/false)"})
				return
			}
			whereClauses = append(whereClauses, param+" = ?")
			queryArgs = append(queryArgs, val)
		}
	}

	// year parameter
	if yearParam := c.Query("jahr"); yearParam != "" {
		year, err := strconv.Atoi(yearParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, HTTPError{Error: "invalid year format"})
			return
		}
		whereClauses = append(whereClauses, "jahr = ?")
		queryArgs = append(queryArgs, year)
	}

	// month parameter
	if monthParam := c.Query("monat"); monthParam != "" {
		month, err := strconv.Atoi(monthParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, HTTPError{Error: "invalid month format"})
			return
		}
		whereClauses = append(whereClauses, "monat = ?")
		queryArgs = append(queryArgs, month)
	}

	// bounding box parameters (latitude & longitude)
	geoParams := []struct {
		param string
		col   string
		op    string
	}{
		{"min_lat", "latitude", ">="},
		{"max_lat", "latitude", "<="},
		{"min_lon", "longitude", ">="},
		{"max_lon", "longitude", "<="},
	}

	for _, gp := range geoParams {
		if valStr := c.Query(gp.param); valStr != "" {
			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, HTTPError{Error: "invalid " + gp.param + " format (must be a float)"})
				return
			}
			whereClauses = append(whereClauses, gp.col+" "+gp.op+" ?")
			queryArgs = append(queryArgs, val)
		}
	}

	// construct final query with WHERE clauses
	finalQuery := baseQuery
	if len(whereClauses) > 0 {
		finalQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// execute query
	rows, err := h.DB.Query(finalQuery, queryArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
		return
	}
	defer rows.Close()

	// JSON response
	var results []data.Unfall
	for rows.Next() {
		var u data.Unfall
		err := rows.Scan(
			&u.Bundesland, &u.Regierungsbezirk, &u.Kreis, &u.Gemeinde, &u.Jahr, &u.Monat, &u.Uhrzeit, &u.Wochentag,
			&u.Schweregrad, &u.Unfallart, &u.Unfalltyp, &u.Lichtverhaeltnis, &u.MitFahrrad, &u.MitPKW,
			&u.MitFussgaenger, &u.MitKraftrad, &u.MitGueterkraftfahrzeug, &u.MitSonstigenVerkehrsmittel,
			&u.Strassenzustand, &u.Latitude, &u.Longitude,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
			return
		}
		results = append(results, u)
	}

	c.JSON(http.StatusOK, results)
}
