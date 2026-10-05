package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/api/common"
	"github.com/rangodisco/yhar/api/dto/request"
	"github.com/rangodisco/yhar/api/services"
)

type AuthHandler struct {
	service *services.AuthService
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func NewAuthHandler(a *services.AuthService) *AuthHandler {
	return &AuthHandler{service: a}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body request.LoginRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	token, refreshToken, err := h.service.HandleUserLogin(c.Request.Context(), body)
	if err != nil {
		common.RespondWithError(c, http.StatusUnauthorized, err, "Invalid credentials")
		return
	}

	common.RespondWithData(c, http.StatusOK, &LoginResponse{AccessToken: token, RefreshToken: refreshToken})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var body request.RefreshRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	token, err := h.service.RefreshToken(body.RefreshToken)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to create new token")
		return
	}

	common.RespondWithData(c, http.StatusOK, &RefreshResponse{AccessToken: token})
}
