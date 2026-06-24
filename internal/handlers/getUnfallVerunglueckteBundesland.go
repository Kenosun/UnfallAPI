package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/gin-gonic/gin"
)

// GetUnfallVerunglueckteBundesland godoc
//
// @Summary      Verunglückte nach Bundesland abrufen
// @Description  (46241-0023/46241-0024) Verunglückte: Bundesländer, Jahre, Monate, Ortslage, Schwere der Verletzung
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Accept       json
// @Produce      json
// @Param        bundesland   query    string  false  "Filter nach Bundesland" enums(Baden-Württemberg, Bayern, Berlin, Brandenburg, Bremen, Hamburg, Hessen, Mecklenburg-Vorpommern, Niedersachsen, Nordrhein-Westfalen, Rheinland-Pfalz, Saarland, Sachsen, Sachsen-Anhalt, Schleswig-Holstein, Thüringen)
// @Param        ortslage     query    string  false  "Filter nach Ortslage (innerorts, außerorts (ohne Autobahnen), auf Autobahnen, Insgesamt)"
// @Param        schweregrad  query    string  false  "Filter nach Schweregrad" enums(Getötete, Schwerverletzte, Leichtverletzte, Insgesamt)
// @Param        jahr         query    int     false  "Filter nach Jahr"
// @Param        monat        query    int     false  "Filter nach Monat (1-12, 0 für Ganzjahresdaten)" minimum(0) maximum(12)
// @Success      200          {array}  data.UnfallVerunglueckteBundesland
// @Failure      400          {object} HTTPError "Bad Request - Invalid parameter type or range"
// @Failure      500          {object} HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallVerunglueckteBundesland [get]
func (h *AccidentHandler) GetUnfallVerunglueckteBundesland(c *gin.Context) {
	baseQuery := `SELECT bundesland, ortslage, schweregrad, jahr, monat, anzahl FROM unfall_verunglueckte_bundesland`

	// dynamic WHERE clauses
	var whereClauses []string
	var queryArgs []any

	// string parameters
	stringParams := map[string]string{
		"bundesland":  c.Query("bundesland"),
		"ortslage":    c.Query("ortslage"),
		"schweregrad": c.Query("schweregrad"),
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
	var results []data.UnfallVerunglueckteBundesland
	for rows.Next() {
		var s data.UnfallVerunglueckteBundesland
		err := rows.Scan(&s.Bundesland, &s.Ortslage, &s.Schweregrad, &s.Jahr, &s.Monat, &s.Anzahl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
			return
		}
		results = append(results, s)
	}

	c.JSON(http.StatusOK, results)
}

// GetUnfallVerunglueckteBundeslandJahre godoc
//
// @Summary      Verfügbare Jahre abrufen
// @Description  Gibt alle Jahre zurück, für die Daten vorhanden sind.
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Produce      json
// @Success      200         {object}  YearsResponse
// @Failure      500         {object}  HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallVerunglueckteBundesland/jahre [get]
func (h *AccidentHandler) GetUnfallVerunglueckteBundeslandJahre(c *gin.Context) {
	query := `SELECT DISTINCT jahr FROM unfall_verunglueckte_bundesland ORDER BY jahr`
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
