package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/gin-gonic/gin"
)

// GetUnfallStrassenverkehrBundesland godoc
//
// @Summary      Straßenverkehrsunfälle mit Personenschaden nach Bundesland abrufen
// @Description  (46241-0022) Straßenverkehrsunfälle mit Personenschaden: Bundesländer, Jahre, Straßenklasse, Ortslage
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Accept       json
// @Produce      json
// @Param        bundesland     query    string  false  "Filter nach Bundesland" enums(Baden-Württemberg, Bayern, Berlin, Brandenburg, Bremen, Hamburg, Hessen, Mecklenburg-Vorpommern, Niedersachsen, Nordrhein-Westfalen, Rheinland-Pfalz, Saarland, Sachsen, Sachsen-Anhalt, Schleswig-Holstein, Thüringen)
// @Param        strassenklasse query    string  false  "Filter nach Straßenklasse" enums(Autobahnen, Bundesstraßen, Landesstraßen, Kreisstraßen, Andere Straßen, Insgesamt)
// @Param        ortslage       query    string  false  "Filter nach Ortslage" enums(innerorts, außerorts, Insgesamt)
// @Param        jahr           query    int     false  "Filter nach Jahr"
// @Success      200            {array}  data.UnfallStrassenverkehrBundesland
// @Failure      400           	{object} HTTPError "Bad Request - Invalid parameter type or range"
// @Failure      500           	{object} HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallStrassenverkehrBundesland [get]
func (h *AccidentHandler) GetUnfallStrassenverkehrBundesland(c *gin.Context) {
	baseQuery := `SELECT bundesland, strassenklasse, ortslage, jahr, anzahl FROM unfall_strassenverkehr_bundesland`

	// dynamic WHERE clauses
	var whereClauses []string
	var queryArgs []any

	// string parameters
	stringParams := map[string]string{
		"bundesland":     c.Query("bundesland"),
		"strassenklasse": c.Query("strassenklasse"),
		"ortslage":       c.Query("ortslage"),
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
	var results []data.UnfallStrassenverkehrBundesland
	for rows.Next() {
		var s data.UnfallStrassenverkehrBundesland
		err := rows.Scan(&s.Bundesland, &s.Strassenklasse, &s.Ortslage, &s.Jahr, &s.Anzahl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
			return
		}
		results = append(results, s)
	}

	c.JSON(http.StatusOK, results)
}

// GetUnfallStrassenverkehrBundeslandJahre godoc
//
// @Summary      Verfügbare Jahre abrufen
// @Description  Gibt alle Jahre zurück, für die Daten vorhanden sind.
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Produce      json
// @Success      200         {object}  YearsResponse
// @Failure      500         {object}  HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallStrassenverkehrBundesland/jahre [get]
func (h *AccidentHandler) GetUnfallStrassenverkehrBundeslandJahre(c *gin.Context) {
	query := `SELECT DISTINCT jahr FROM unfall_strassenverkehr_bundesland ORDER BY jahr`
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
