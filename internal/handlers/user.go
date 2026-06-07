package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// GetUser godoc
//
// @Summary Get user by ID
// @Description Returns a single user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")

	if id == "0" {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, User{
		ID:   id,
		Name: "John Doe",
	})
}
