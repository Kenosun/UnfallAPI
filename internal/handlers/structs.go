package handlers

import "database/sql"

type AccidentHandler struct {
	DB *sql.DB
}

type HTTPError struct {
	Error string `json:"error" example:"invalid parameter format"`
}

type JahreResponse struct {
	Jahre []int `json:"jahre" example:"[2016,2017,2018,2019,2020,2021,2022,2023,2024]"`
}
