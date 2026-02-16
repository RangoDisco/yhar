package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/services"
)

type UserHandler struct {
	authService *services.AuthService
}

type GetUserResponse struct {
	User UserWithID `json:"user"`
}

type UserWithID struct {
	ID int64 `json:"id"`
}

func NewUserHandler(a *services.AuthService) *UserHandler {
	return &UserHandler{authService: a}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	queriedID := c.Param("userID")
	// TODO implement
	if queriedID != "me" {
		common.RespondWithError(c, http.StatusNotFound, errors.New("route not implemented"), "User not found")
		return
	}

	rawUser, exists := c.Get("user")
	if !exists {
		common.RespondWithError(c, http.StatusUnauthorized, errors.New("context's user not found"), "Unauthorized")
		return
	}

	user, ok := rawUser.(*dto.UserPassport)
	if !ok {
		common.RespondWithError(c, http.StatusUnauthorized, errors.New("unable to convert context's to model"), "Unauthorized")
		return
	}

	common.RespondWithData(c, http.StatusOK, &GetUserResponse{User: UserWithID{user.ID}})
}
