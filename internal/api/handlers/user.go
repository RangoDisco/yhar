package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/services"
)

type UserHandler struct {
	authService *services.AuthService
}

func NewUserHandler(a *services.AuthService) *UserHandler {
	return &UserHandler{authService: a}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	queriedID := c.Param("userID")
	// TODO implement
	if queriedID != "me" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	rawUser, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := rawUser.(*models.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"user": gin.H{"id": user.ID}}})
}
