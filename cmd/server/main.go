package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Kenosun/UnfallAPI/docs"

	"github.com/Kenosun/UnfallAPI/internal/handlers"
)

// @title UnfallAPI
// @version 1.0
// @description REST API with Swagger documentation.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	handlers.LoadUnfallData()

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/users/:id", handlers.GetUser)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
