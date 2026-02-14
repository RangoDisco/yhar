package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/types/auth"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(a *services.AuthService) *AuthHandler {
	return &AuthHandler{service: a}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body auth.LoginRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.HandleUserLogin(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"token": token},
	})
}
