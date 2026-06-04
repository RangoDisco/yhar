package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/tests/factories"
	"github.com/stretchr/testify/assert"
)

// TODO handle admin case for all
// TODO handle anon case

// TestGetTopArtists tests cases when fetching a user's top artists
func TestGetTopArtists(t *testing.T) {
	db := SetupDB(t)
	privateUser, regularUser := factories.GetStatsSeedData(t, db)

	t.Run("Private user accessing their own artist data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/artists", privateUser.ID), nil)
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

		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/artists", privateUser.ID), nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Regular user accessing public user's artist data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/artists", regularUser.ID), nil)
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
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/albums", privateUser.ID), nil)
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

		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("api/users/%d/scrobbles/top/albums", privateUser.ID), nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Regular user accessing public user's album data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/albums", regularUser.ID), nil)
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
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/tracks", privateUser.ID), nil)
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

		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/tracks", privateUser.ID), nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Regular user accessing public user's track data", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/tracks", regularUser.ID), nil)
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
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/history", privateUser.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.HistoryResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 8)
		assert.NotNil(t, result.Data.Results[0].Track.Title)
		assert.NotNil(t, result.Data.Results[0].Track.Album)
		assert.NotEmpty(t, result.Data.Results[0].Track.Artists)
		assert.NotNil(t, result.Data.Results[0].Track.Artists[0].Name)
		assert.NotNil(t, result.Data.Results[0].ScrobbledAt)
	})

	t.Run("Regular user accessing another private user's history", func(t *testing.T) {
		router := SetupRouter(t, db, &regularUser)

		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/top/tracks", privateUser.ID), nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, 1, 1)
	})

	t.Run("Regular user accessing public user's history", func(t *testing.T) {
		router := SetupRouter(t, db, &privateUser)
		w := doRequest(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d/scrobbles/history", regularUser.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var result common.APIResponse[common.PaginatedResponse[[]dto.HistoryResult]]
		body, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(body, &result)

		assert.Len(t, result.Data.Results, 8)
		assert.NotNil(t, result.Data.Results[0].Track.Title)
		assert.NotNil(t, result.Data.Results[0].Track.Album)
		assert.NotEmpty(t, result.Data.Results[0].Track.Artists)
		assert.NotNil(t, result.Data.Results[0].Track.Artists[0].Name)
		assert.NotNil(t, result.Data.Results[0].ScrobbledAt)
	})

	t.Cleanup(func() { db.Rollback() })
}
