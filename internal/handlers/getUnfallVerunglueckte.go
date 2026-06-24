package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/gin-gonic/gin"
)

// GetUnfallVerunglueckte godoc
//
// @Summary      Verunglückte abrufen
// @Description  (46241-0007/46241-0008) Verunglückte: Deutschland, Jahre, Monate, Geschlecht, Altersgruppen, Art der Verkehrsbeteiligung, Ortslage, Schwere der Verletzung
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Accept       json
// @Produce      json
// @Param        verkehrsart  query    string  false  "Filter nach Art der Verkehrsbeteiligung" enums(Kraftrad mit Versicherungskennzeichen, Kraftrad mit amtlichem Kennzeichen, Elektrokleinstfahrzeuge, Personenkraftwagen, Kraftomnibus, Güterkraftfahrzeug, Landwirtschaftliche Zugmaschine, Übrige Kraftfahrzeuge, Fahrrad ohne Hilfsmotor, Pedelecs, Andere Fahrzeuge, Fußgänger, Andere Personen, Insgesamt)
// @Param        ortslage     query    string  false  "Filter nach Ortslage (innerorts, außerorts (ohne Autobahnen), auf Autobahnen, Insgesamt)"
// @Param        schweregrad  query    string  false  "Filter nach Schweregrad" enums(Getötete, Schwerverletzte, Leichtverletzte, Insgesamt)
// @Param        geschlecht   query    string  false  "Filter nach Geschlecht" enums(männlich, weiblich, Ohne Angabe, Insgesamt)
// @Param        altersgruppe query    string  false  "Filter nach Altersgruppe" enums(unter 15 Jahre, 15 bis unter 18 Jahre, 18 bis unter 21 Jahre, 21 bis unter 25 Jahre, 25 bis unter 35 Jahre, 35 bis unter 45 Jahre, 45 bis unter 55 Jahre, 55 bis unter 65 Jahre, 65 bis unter 75 Jahre, 75 Jahre und mehr, Alter unbekannt, Insgesamt)"
// @Param        jahr         query    int     false  "Filter nach Jahr"
// @Param        monat        query    int     false  "Filter nach Monat (1-12, 0 für Ganzjahresdaten)" minimum(0) maximum(12)
// @Success      200          {array}  data.UnfallVerunglueckte
// @Failure      400          {object} HTTPError "Bad Request - Invalid parameter type or range"
// @Failure      500          {object} HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallVerunglueckte [get]
func (h *AccidentHandler) GetUnfallVerunglueckte(c *gin.Context) {
	baseQuery := `SELECT verkehrsart, ortslage, schweregrad, geschlecht, altersgruppe, jahr, monat, anzahl FROM unfall_verunglueckte`

	// dynamic WHERE clauses
	var whereClauses []string
	var queryArgs []any

	// string parameters
	stringParams := map[string]string{
		"verkehrsart":  c.Query("verkehrsart"),
		"ortslage":     c.Query("ortslage"),
		"schweregrad":  c.Query("schweregrad"),
		"geschlecht":   c.Query("geschlecht"),
		"altersgruppe": c.Query("altersgruppe"),
	}

	for column, value := range stringParams {
		if value != "" {
			whereClauses = append(whereClauses, column+" = ?")
			queryArgs = append(queryArgs, value)
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
	var results []data.UnfallVerunglueckte
	for rows.Next() {
		var s data.UnfallVerunglueckte
		err := rows.Scan(&s.Verkehrsart, &s.Ortslage, &s.Schweregrad, &s.Geschlecht, &s.Altersgruppe, &s.Jahr, &s.Monat, &s.Anzahl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
			return
		}
		results = append(results, s)
	}

	c.JSON(http.StatusOK, results)
}

// GetUnfallVerunglueckteJahre godoc
//
// @Summary      Verfügbare Jahre abrufen
// @Description  Gibt alle Jahre zurück, für die Daten vorhanden sind.
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Produce      json
// @Success      200         {object}  YearsResponse
// @Failure      500         {object}  HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallVerunglueckte/jahre [get]
func (h *AccidentHandler) GetUnfallVerunglueckteJahre(c *gin.Context) {
	query := `SELECT DISTINCT jahr FROM unfall_verunglueckte ORDER BY jahr`
	rows, err := h.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
		return
	}
	defer rows.Close()

	var years []int
	for rows.Next() {
		var year int
		if err := rows.Scan(&year); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
			return
		}
		years = append(years, year)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, YearsResponse{Years: years})
}
