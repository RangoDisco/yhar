package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/handlers"
	"github.com/rangodisco/yhar/tests/factories"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	db := SetupDB(t)

	privateUser := factories.SeedUser(t, db, "private", "123", "USER", false)
	publicUser := factories.SeedUser(t, db, "public", "123", "USER", true)

	t.Run("User getting its own data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/me"), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[handlers.GetUserResponse]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Equal(t, privateUser.ID, result.Data.User.ID)
	})

	// For future reference
	// TODO: implement get user data by id
	t.Run("User getting an other user data by its ID", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d", publicUser.ID), nil)
		// Route is yet to be implemented
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Cleanup(func() { db.Rollback() })
}
