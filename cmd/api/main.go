package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Kenosun/UnfallAPI/docs"

	"github.com/Kenosun/UnfallAPI/internal/database"
	"github.com/Kenosun/UnfallAPI/internal/handlers"
	"github.com/Kenosun/UnfallAPI/internal/service"
)

// @title 		UnfallAPI
// @version 	1.0
// @description	REST API with Swagger documentation.
// @host		localhost:8080
// @BasePath 	/api/v1
func main() {
	dbName := "unfallData.db"
	port := "8080"

	// configure log options
	log.SetReportTimestamp(true)
	log.SetLevel(log.DebugLevel)

	// check if the database file already exists
	_, err := os.Stat(dbName)
	dbExists := err == nil

	// download unfallData if the database file doesn't exist
	if !dbExists {
		log.Debug("Database file not found. Starting initial setup...")
		if err := service.DownloadUnfallData(); err != nil {
			log.Error("Failed to download UnfallData", "error", err)
		}
	} else {
		log.Debug("Database file found. Skipping unfallData download.")
	}

	// initialize database
	db, err := database.InitializeDB("unfallData.db")
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// load unfallData into the database if the database didn't exist previously
	if !dbExists {
		log.Debug("Loading unfallData into database...")
		if err := database.LoadUnfallData(db); err != nil {
			log.Error("Failed to load UnfallData", "error", err)
		}
	} else {
		log.Debug("Database already exists. Skipping unfallData loading.")
	}

	println("")

	r := gin.Default()
	h := handlers.AccidentHandler{DB: db}

	api := r.Group("/api/v1")
	{
		api.GET("/unfallStatistik", h.GetUnfallStatistik)
		api.GET("/unfallStrassenverkehr", h.GetUnfallStrassenverkehr)
		api.GET("/unfallPersonenschaden", h.GetUnfallPersonenschaden)
		api.GET("/unfallVerunglueckte", h.GetUnfallVerunglueckte)
		api.GET("/unfallFehlverhalten", h.GetUnfallFehlverhalten)
		api.GET("/unfallBeteiligung", h.GetUnfallBeteiligung)
		api.GET("/unfallStatistikBundesland", h.GetUnfallStatistikBundesland)
		api.GET("/unfallStrassenverkehrBundesland", h.GetUnfallStrassenverkehrBundesland)
		api.GET("/unfallVerunglueckteBundesland", h.GetUnfallVerunglueckteBundesland)
		api.GET("/unfall", h.GetUnfall)
		api.GET("/ort", h.GetOrt)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	println("")
	log.Info("API endpoints available", "url", "http://localhost:"+port+"/api/v1/<endpoint>")
	log.Info("Swagger documentation available", "url", "http://localhost:"+port+"/swagger/index.html")
	println("")

	r.Run(":" + port)
}
