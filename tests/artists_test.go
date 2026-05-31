package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/tests/factories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatchArtist(t *testing.T) {
	db := SetupDB(t)
	regularUser := factories.SeedUser(t, db, "regular", "123", "USER", false)
	adminUser := factories.SeedUser(t, db, "admin", "123", "ADMIN", false)

	tracks := factories.CreateScrobbleContent(t, db)
	body := dto.UpdateArtistInput{Name: new("Updated name")}
	out, err := json.Marshal(body)
	require.NoError(t, err)
	bytesBody := bytes.NewBuffer(out)

	t.Run("Regular user patching an artist", func(t *testing.T) {
		router := SetupRouter(t, db, &regularUser)

		w := doRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/artists/%d", tracks[0].Artists[0].ID), bytesBody)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Admin user patching an artist", func(t *testing.T) {
		router := SetupRouter(t, db, &adminUser)

		w := doRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/artists/%d", tracks[0].Artists[0].ID), bytesBody)
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[*models.Artist]
		resBody, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(resBody, &result)

		assert.Equal(t, "Updated name", result.Data.Name)

		// Double check that the data has been updated in db
		var a models.Artist
		err = db.Where("id = ?", tracks[0].Artists[0].ID).First(&a).Error
		assert.NoError(t, err)

		assert.Equal(t, "Updated name", result.Data.Name)
	})

	t.Cleanup(func() { db.Rollback() })
}
