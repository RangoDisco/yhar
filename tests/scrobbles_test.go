package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/tests/factories"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteScrobble(t *testing.T) {
	db := SetupDB(t)
	privateUser := factories.SeedUser(t, db, "private", "123", "USER", false)
	regularUser := factories.SeedUser(t, db, "regular", "123", "USER", true)
	tracks := factories.CreateScrobbleContent(t, db)

	scrobbles := factories.CreateScrobbles(t, db, tracks[0].ID, privateUser.ID, 1)

	t.Run("Regular user deleting another user's scrobble", func(t *testing.T) {
		router := SetupRouter(t, db, &regularUser)
		w := doRequest(t, router, http.MethodDelete, fmt.Sprintf("/api/users/%d/scrobbles/%d", privateUser.ID, scrobbles[0].ID), nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Private user deleting their own scrobble", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, http.MethodDelete, fmt.Sprintf("/api/scrobbles/%d", scrobbles[0].ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		// Double check that the scrobble has been deleted in db
		var s models.Scrobble
		err := db.Where("id = ?", scrobbles[0].ID).First(&s).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Cleanup(func() { db.Rollback() })

}
