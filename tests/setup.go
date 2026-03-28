package tests

import (
	"context"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/config"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("DSN")
	db, _ := gorm.Open(postgres.Open(dsn))

	tx := db.Begin()

	alTypeRes := gorm.WithResult()
	err := gorm.G[any](db, alTypeRes).Exec(context.Background(), "CREATE type album_type AS enum ('ALBUM', 'SINGLE', 'EP', 'COMPILATION');")
	require.NoError(t, err)

	scTypRes := gorm.WithResult()
	err = gorm.G[any](db, scTypRes).Exec(context.Background(), "CREATE type scrobble_origin AS enum ('SUBSONIC');")
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.Album{},
		&models.Artist{},
		&models.Genre{},
		&models.Image{},
		&models.Scrobble{},
		&models.Track{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Session{},
	)
	require.NoError(t, err)

	var roles = []*models.Role{{Name: "ADMIN"}, {Name: "USER"}}
	db.WithContext(context.Background()).CreateInBatches(&roles, 2)
	require.NoError(t, err)

	t.Cleanup(func() { tx.Rollback() })

	return tx
}

func SetupRouter(t *testing.T, db *gorm.DB) (*gin.Engine, *config.Handlers) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	_, _, handlers, _, _ := config.AutoWire(db)

	// Mock authentication middleware, TODO: remove
	r.Use(func(c *gin.Context) {
		c.Set("user", &dto.UserPassport{ID: 1, Username: "rango"})
		c.Next()
	})

	return r, handlers
}
