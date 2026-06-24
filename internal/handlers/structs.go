package handlers

import "database/sql"

type AccidentHandler struct {
	DB *sql.DB
}

type HTTPError struct {
	Error string `json:"error" example:"invalid parameter format"`
}

type YearsResponse struct {
	Years []int `json:"jahre" example:"2022"`
}
