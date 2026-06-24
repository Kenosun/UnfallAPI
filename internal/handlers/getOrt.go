package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kenosun/UnfallAPI/internal/data"
	"github.com/gin-gonic/gin"
)

// GetOrt godoc
//
// @Summary      Detaillierte Ortsdaten im Gemeindeverzeichnis abrufen
// @Description  Strukturierte Daten über Bundes­länder, Regierungs­bezirke, Kreise, Gemeindeverbände und Gemeinden abrufen inklusive Bevölkerung, Fläche, und Geokoordinaten.
// @Tags         Gemeindeverzeichnis des Statistischen Bundesamtes
// @Accept       json
// @Produce      json
// @Param        bundesland     		query    string   false  "Filter nach Bundesland" enums(Baden-Württemberg, Bayern, Berlin, Brandenburg, Bremen, Hamburg, Hessen, Mecklenburg-Vorpommern, Niedersachsen, Nordrhein-Westfalen, Rheinland-Pfalz, Saarland, Sachsen, Sachsen-Anhalt, Schleswig-Holstein, Thüringen)
// @Param        regierungsbezirk 		query    string   false  "Filter nach Regierungsbezirk (Amtlicher Gemeindeschlüssel)"
// @Param        kreis          		query    string   false  "Filter nach Kreis (Amtlicher Gemeindeschlüssel)"
// @Param        gemeinde       		query    string   false  "Filter nach Gemeinde (Amtlicher Gemeindeschlüssel)"
// @Param        name           		query    string   false  "Filter nach Ortsname"
// @Param        gemeindeverband 		query    string   false  "Filter nach Gemeindeverband"
// @Param        landkreis      		query    string   false  "Filter nach Landkreis"
// @Param        postleitzahl   		query    string   false  "Filter nach Postleitzahl"
// @Param        min_flaeche            query    float64  false  "Filter nach Mindestfläche"
// @Param        max_flaeche            query    float64  false  "Filter nach Maximalfläche"
// @Param        min_bevoelkerung       query    int      false  "Filter nach Mindestbevölkerung"
// @Param        max_bevoelkerung       query    int      false  "Filter nach Maximalbevölkerung"
// @Param        min_maennlich          query    int      false  "Filter nach Mindestanzahl männlich"
// @Param        max_maennlich          query    int      false  "Filter nach Maximalanzahl männlich"
// @Param        min_weiblich           query    int      false  "Filter nach Mindestanzahl weiblich"
// @Param        max_weiblich           query    int      false  "Filter nach Maximalanzahl weiblich"
// @Param        reisegebiet    		query    string   false  "Filter nach Reisegebiet"
// @Param        verstaedterungsgrad 	query    string   false  "Filter nach Verstädterungsgrad" enums(dicht besiedelt, mittlere Besiedlungsdichte, gering besiedelt)
// @Param        min_lat          		query    float64  false  "Minimum Latitude"
// @Param        max_lat          		query    float64  false  "Maximum Latitude"
// @Param        min_lon          		query    float64  false  "Minimum Longitude"
// @Param        max_lon          		query    float64  false  "Maximum Longitude"
// @Success      200            		{array}  data.Ort
// @Failure      400         			{object} HTTPError "Bad Request - Invalid parameter type or range"
// @Failure      500         			{object} HTTPError "Internal Server Error - Database execution or scanning failure"
// @Router       /ort [get]
func (h *AccidentHandler) GetOrt(c *gin.Context) {
	baseQuery := `SELECT bundesland, regierungsbezirk, kreis, gemeinde, name, gemeindeverband, landkreis, postleitzahl, flaeche, bevoelkerung, maennlich, weiblich, reisegebiet, verstaedterungsgrad, latitude, longitude FROM ort`

	// dynamic WHERE clauses
	var whereClauses []string
	var queryArgs []any

	// string parameters
	stringParams := map[string]string{
		"bundesland":          c.Query("bundesland"),
		"regierungsbezirk":    c.Query("regierungsbezirk"),
		"kreis":               c.Query("kreis"),
		"gemeinde":            c.Query("gemeinde"),
		"name":                c.Query("name"),
		"gemeindeverband":     c.Query("gemeindeverband"),
		"landkreis":           c.Query("landkreis"),
		"postleitzahl":        c.Query("postleitzahl"),
		"reisegebiet":         c.Query("reisegebiet"),
		"verstaedterungsgrad": c.Query("verstaedterungsgrad"),
	}

	for column, value := range stringParams {
		if value != "" {
			whereClauses = append(whereClauses, column+" = ?")
			queryArgs = append(queryArgs, value)
		}
	}

	// numeric range parameters
	numParams := []struct {
		param   string
		col     string
		op      string
		isFloat bool
	}{
		{"min_flaeche", "flaeche", ">=", true},
		{"max_flaeche", "flaeche", "<=", true},
		{"min_bevoelkerung", "bevoelkerung", ">=", false},
		{"max_bevoelkerung", "bevoelkerung", "<=", false},
		{"min_maennlich", "maennlich", ">=", false},
		{"max_maennlich", "maennlich", "<=", false},
		{"min_weiblich", "weiblich", ">=", false},
		{"max_weiblich", "weiblich", "<=", false},
	}

	for _, np := range numParams {
		if valStr := c.Query(np.param); valStr != "" {
			var val any
			var err error

			if np.isFloat {
				val, err = strconv.ParseFloat(valStr, 64)
			} else {
				val, err = strconv.Atoi(valStr)
			}

			if err != nil {
				c.JSON(http.StatusBadRequest, HTTPError{Error: "invalid " + np.param + " format (must be a number)"})
				return
			}

			whereClauses = append(whereClauses, np.col+" "+np.op+" ?")
			queryArgs = append(queryArgs, val)
		}
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
	var results []data.Ort
	for rows.Next() {
		var o data.Ort
		err := rows.Scan(
			&o.Bundesland,
			&o.Regierungsbezirk,
			&o.Kreis,
			&o.Gemeinde,
			&o.Name,
			&o.Gemeindeverband,
			&o.Landkreis,
			&o.Postleitzahl,
			&o.Flaeche,
			&o.Bevoelkerung,
			&o.Maennlich,
			&o.Weiblich,
			&o.Reisegebiet,
			&o.Verstaedterungsgrad,
			&o.Latitude,
			&o.Longitude,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
			return
		}
		results = append(results, o)
	}

	c.JSON(http.StatusOK, results)
}
