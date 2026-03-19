package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/types/auth"
)

type AuthHandler struct {
	service *services.AuthService
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(a *services.AuthService) *AuthHandler {
	return &AuthHandler{service: a}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body auth.LoginRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	token, err := h.service.HandleUserLogin(c.Request.Context(), body)
	if err != nil {
		common.RespondWithError(c, http.StatusUnauthorized, err, "Invalid credentials")
		return
	}

	common.RespondWithData(c, http.StatusOK, &LoginResponse{Token: token})
}
