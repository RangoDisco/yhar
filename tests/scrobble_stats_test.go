package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/tests/factories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO handle admin case for all

// TestGetTopArtists tests cases when fetching a user's top artists
func TestGetTopArtists(t *testing.T) {
	db := SetupDB(t)
	privateUser, regularUser := factories.GetStatsSeedData(t, db)

	t.Run("Private user accessing their own artist data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, privateUser.ID, "/scrobbles/top/artists")
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.TopArtistResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 3)
		assert.NotNil(t, result.Data.Results[0].Name)
		assert.Equal(t, *result.Data.Results[0].ScrobbleCount, 3)
	})

	t.Run("Regular user accessing another private user's artist data", func(t *testing.T) {
		router := SetupRouter(t, db, &regularUser)

		w := doRequest(t, router, privateUser.ID, "/scrobbles/top/artists")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Regular user accessing public user's artist data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, regularUser.ID, "/scrobbles/top/artists")
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.TopArtistResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 3)
		assert.Equal(t, *result.Data.Results[0].ScrobbleCount, 3)
	})

	t.Cleanup(func() { db.Rollback() })
}

// TestGetTopAlbums tests cases when fetching a user's top albums
func TestGetTopAlbums(t *testing.T) {
	db := SetupDB(t)
	privateUser, regularUser := factories.GetStatsSeedData(t, db)

	t.Run("Private user accessing their own album data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, privateUser.ID, "/scrobbles/top/albums")
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.TopAlbumResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 3)
		assert.NotNil(t, result.Data.Results[0].Title)
		assert.Equal(t, *result.Data.Results[0].ScrobbleCount, 3)
	})

	t.Run("Regular user accessing another private user's album data", func(t *testing.T) {
		router := SetupRouter(t, db, &regularUser)

		w := doRequest(t, router, privateUser.ID, "/scrobbles/top/albums")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Regular user accessing public user's album data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, regularUser.ID, "/scrobbles/top/albums")
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.TopAlbumResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 3)
		assert.NotNil(t, result.Data.Results[0].Title)
		assert.Equal(t, *result.Data.Results[0].ScrobbleCount, 3)
	})

	t.Cleanup(func() { db.Rollback() })
}

// TestGetTopTracks tests cases when fetching a user's top tracks
func TestGetTopTracks(t *testing.T) {
	db := SetupDB(t)
	privateUser, regularUser := factories.GetStatsSeedData(t, db)

	t.Run("Private user accessing their own track data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, privateUser.ID, "/scrobbles/top/tracks")
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.TrackResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 8)
		assert.NotNil(t, result.Data.Results[0].Title)
		assert.NotNil(t, result.Data.Results[0].Album)
		assert.NotEmpty(t, result.Data.Results[0].Artists)
		assert.NotNil(t, result.Data.Results[0].Artists[0].Name)
		assert.Equal(t, 1, *result.Data.Results[0].ScrobbleCount)
	})

	t.Run("Regular user accessing another private user's track data", func(t *testing.T) {
		router := SetupRouter(t, db, &regularUser)

		w := doRequest(t, router, privateUser.ID, "/scrobbles/top/tracks")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Regular user accessing public user's track data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, regularUser.ID, "/scrobbles/top/tracks")
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.TrackResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 8)
		assert.NotNil(t, result.Data.Results[0].Title)
		assert.NotNil(t, result.Data.Results[0].Album)
		assert.NotEmpty(t, result.Data.Results[0].Artists)
		assert.NotNil(t, result.Data.Results[0].Artists[0].Name)
		assert.Equal(t, 1, *result.Data.Results[0].ScrobbleCount)
	})

	t.Cleanup(func() { db.Rollback() })
}

// TestGetHistory tests cases when fetching a user's listening history
func TestGetHistory(t *testing.T) {
	db := SetupDB(t)
	privateUser, regularUser := factories.GetStatsSeedData(t, db)

	t.Run("Private user accessing their own history", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, privateUser.ID, "/scrobbles/history")
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.TrackResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 8)
		assert.NotNil(t, result.Data.Results[0].Title)
		assert.NotNil(t, result.Data.Results[0].Album)
		assert.NotEmpty(t, result.Data.Results[0].Artists)
		assert.NotNil(t, result.Data.Results[0].Artists[0].Name)
		assert.Equal(t, 1, *result.Data.Results[0].ScrobbleCount)
	})

	t.Run("Regular user accessing another private user's history", func(t *testing.T) {
		router := SetupRouter(t, db, &regularUser)

		w := doRequest(t, router, privateUser.ID, "/scrobbles/top/tracks")
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, 1, 1)
	})

	t.Run("Regular user accessing public user's history", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, regularUser.ID, "/scrobbles/history")
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.TrackResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 8)
		assert.NotNil(t, result.Data.Results[0].Title)
		assert.NotNil(t, result.Data.Results[0].Album)
		assert.NotEmpty(t, result.Data.Results[0].Artists)
		assert.NotNil(t, result.Data.Results[0].Artists[0].Name)
		assert.Equal(t, 1, *result.Data.Results[0].ScrobbleCount)
	})
}

func doRequest(t *testing.T, r *gin.Engine, uID int64, path string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/api/users/%d%s", uID, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
