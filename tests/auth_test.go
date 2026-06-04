package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/handlers"
	"github.com/rangodisco/yhar/internal/api/types/auth"
	"github.com/rangodisco/yhar/tests/factories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	db := SetupDB(t)
	_ = factories.SeedUser(t, db, "regularUser", "123", "USER", false)

	t.Run("Login with correct credentials", func(t *testing.T) {
		body := auth.LoginRequest{Username: "regularUser", Password: "123"}
		out, err := json.Marshal(body)
		require.NoError(t, err)
		bytesBody := bytes.NewBuffer(out)

		router := SetupRouter(t, db, nil)
		w := doRequest(t, router, http.MethodPost, "/api/auth/login", bytesBody)
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[handlers.LoginResponse]
		resBody, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(resBody, &result)

		assert.NotEmpty(t, result.Data.AccessToken)
		assert.NotEmpty(t, result.Data.RefreshToken)
	})

	t.Run("Login with incorrect credentials", func(t *testing.T) {
		body := auth.LoginRequest{Username: "regularUser", Password: "12345"}
		out, err := json.Marshal(body)
		require.NoError(t, err)
		bytesBody := bytes.NewBuffer(out)

		router := SetupRouter(t, db, nil)
		w := doRequest(t, router, http.MethodPost, "/api/auth/login", bytesBody)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Cleanup(func() { db.Rollback() })
}

func TestRefresh(t *testing.T) {
	db := SetupDB(t)
	regularUser := factories.SeedUser(t, db, "regularUser", "123", "USER", false)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": regularUser.Username,
		"role":     regularUser.Role.Name,
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	assert.NoError(t, err)

	t.Run("Refresh token with valid refresh token", func(t *testing.T) {
		router := SetupRouter(t, db, nil)

		body := auth.RefreshRequest{RefreshToken: refreshTokenString}
		out, err := json.Marshal(body)
		require.NoError(t, err)
		bytesBody := bytes.NewBuffer(out)

		w := doRequest(t, router, http.MethodPost, "/api/auth/refresh", bytesBody)
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[handlers.RefreshResponse]
		resBody, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(resBody, &result)
		assert.NotEmpty(t, result.Data.AccessToken)
	})

	t.Run("Refresh token with invalid refresh token", func(t *testing.T) {
		router := SetupRouter(t, db, nil)

		body := auth.RefreshRequest{RefreshToken: "INVALID"}
		out, err := json.Marshal(body)
		require.NoError(t, err)
		bytesBody := bytes.NewBuffer(out)

		w := doRequest(t, router, http.MethodPost, "/api/auth/refresh", bytesBody)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Cleanup(func() { db.Rollback() })
}
