package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/gin-gonic/gin"
)

// GetUnfallStatistikBundesland godoc
//
// @Summary      Unfälle (polizeilich erfasste) nach Bundesland abrufen
// @Description  (46241-0020/46241-0021) Unfälle (polizeilich erfasste): Bundesländer, Jahre, Monate, Unfallkategorie, Ortslage
// @Tags         GENESIS-Online (Die Datenbank des Statistischen Bundesamtes)
// @Accept       json
// @Produce      json
// @Param        bundesland       query    string  false  "Filter nach Bundesland" enums(Baden-Württemberg, Bayern, Berlin, Brandenburg, Bremen, Hamburg, Hessen, Mecklenburg-Vorpommern, Niedersachsen, Nordrhein-Westfalen, Rheinland-Pfalz, Saarland, Sachsen, Sachsen-Anhalt, Schleswig-Holstein, Thüringen)
// @Param        unfallkategorie  query    string  false  "Filter nach Unfallkategorie" enums(Unfälle mit Personenschaden, Schwerwiegende Unfälle mit Sachschaden i.e.S, Sonst. Unfälle unter dem Einfluss berausch. Mittel, Übrige Sachschadensunfälle, Insgesamt)
// @Param        ortslage         query    string  false  "Filter nach Ortslage (innerorts, außerorts (ohne Autobahnen), auf Autobahnen, Insgesamt)"
// @Param        jahr             query    int     false  "Filter nach Jahr"
// @Param        monat            query    int     false  "Filter nach Monat (1-12, 0 für Ganzjahresdaten)" minimum(0) maximum(12)
// @Success      200              {array}  data.UnfallStatistikBundesland
// @Failure      400              {object} HTTPError "Bad Request - Invalid parameter type or range"
// @Failure      500              {object} HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /unfallStatistikBundesland [get]
func (h *AccidentHandler) GetUnfallStatistikBundesland(c *gin.Context) {
	baseQuery := `SELECT bundesland, unfallkategorie, ortslage, jahr, monat, anzahl FROM unfall_statistik_bundesland`

	// dynamic WHERE clauses
	var whereClauses []string
	var queryArgs []any

	// string parameters
	stringParams := map[string]string{
		"bundesland":      c.Query("bundesland"),
		"unfallkategorie": c.Query("unfallkategorie"),
		"ortslage":        c.Query("ortslage"),
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
	var results []data.UnfallStatistikBundesland
	for rows.Next() {
		var s data.UnfallStatistikBundesland
		err := rows.Scan(&s.Bundesland, &s.Unfallkategorie, &s.Ortslage, &s.Jahr, &s.Monat, &s.Anzahl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
			return
		}
		results = append(results, s)
	}

	c.JSON(http.StatusOK, results)
}
