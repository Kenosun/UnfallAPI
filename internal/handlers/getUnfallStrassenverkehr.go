package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/gin-gonic/gin"
)

// GetUnfallStrassenverkehr godoc
//
// @Summary      Straßenverkehrsunfälle mit Personenschaden abrufen
// @Description  (46241-0003/46241-0004) Straßenverkehrsunfälle mit Personenschaden, Getöteten, Schwer- und Leichtverletzten: Deutschland, Jahre, Monate, Straßenklasse, Ortslage
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Accept       json
// @Produce      json
// @Param        strassenklasse query    string  false  "Filter nach Straßenklasse" enums(Autobahnen, Bundesstraßen, Landesstraßen, Kreisstraßen, Andere Straßen, Insgesamt)
// @Param        ortslage       query    string  false  "Filter nach Ortslage" enums(innerorts, außerorts, Insgesamt)
// @Param        kategorie      query    string  false  "Filter nach Kategorie" enums(Unfälle mit Personenschaden, Getötete, Schwerverletzte, Leichtverletzte)
// @Param        jahr           query    int     false  "Filter nach Jahr"
// @Param        monat          query    int     false  "Filter nach Monat (1-12, 0 für Ganzjahresdaten)" minimum(0) maximum(12)
// @Success      200            {array}  data.UnfallStrassenverkehr
// @Failure      400           	{object} HTTPError "Bad Request - Invalid parameter type or range"
// @Failure      500           	{object} HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallStrassenverkehr [get]
func (h *AccidentHandler) GetUnfallStrassenverkehr(c *gin.Context) {
	baseQuery := `SELECT strassenklasse, ortslage, kategorie, jahr, monat, anzahl FROM unfall_strassenverkehr`

	// dynamic WHERE clauses
	var whereClauses []string
	var queryArgs []any

	// string parameters
	stringParams := map[string]string{
		"strassenklasse": c.Query("strassenklasse"),
		"ortslage":       c.Query("ortslage"),
		"kategorie":      c.Query("kategorie"),
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
	var results []data.UnfallStrassenverkehr
	for rows.Next() {
		var s data.UnfallStrassenverkehr
		err := rows.Scan(&s.Strassenklasse, &s.Ortslage, &s.Kategorie, &s.Jahr, &s.Monat, &s.Anzahl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
			return
		}
		results = append(results, s)
	}

	c.JSON(http.StatusOK, results)
}

// GetUnfallStrassenverkehrJahre godoc
//
// @Summary      Verfügbare Jahre abrufen
// @Description  Gibt alle Jahre zurück, für die Daten vorhanden sind.
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Produce      json
// @Success      200         {object}  YearsResponse
// @Failure      500         {object}  HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallStrassenverkehr/jahre [get]
func (h *AccidentHandler) GetUnfallStrassenverkehrJahre(c *gin.Context) {
	query := `SELECT DISTINCT jahr FROM unfall_strassenverkehr ORDER BY jahr`
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
