package tests

import (
	"context"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/config"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	dsn := os.Getenv("DSN")
	db, _ := gorm.Open(postgres.Open(dsn))

	alTypeRes := gorm.WithResult()
	err := gorm.G[any](db, alTypeRes).Exec(context.Background(), "CREATE type album_type AS enum ('ALBUM', 'SINGLE', 'EP', 'COMPILATION');")
	if err != nil {
		panic(err)
	}

	scTypRes := gorm.WithResult()
	err = gorm.G[any](db, scTypRes).Exec(context.Background(), "CREATE type scrobble_origin AS enum ('SUBSONIC');")
	if err != nil {
		panic(err)
	}

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
	if err != nil {
		panic(err)
	}

	var roles = []*models.Role{{Name: "ADMIN", Permissions: []models.Permission{
		{Name: "UPDATE_ARTIST"},
		{Name: "UPDATE_ALBUM"},
		{Name: "IMAGE_UPLOAD"},
		{Name: "MANUAL_SCROBBLE"},
	}}, {Name: "USER"}}
	err = db.WithContext(context.Background()).CreateInBatches(&roles, 2).Error
	if err != nil {
		panic(err)
	}
}

func SetupDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("DSN")
	db, _ := gorm.Open(postgres.Open(dsn))

	tx := db.Begin()

	return tx
}

func SetupRouter(t *testing.T, db *gorm.DB, caller *models.User) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	_, s, h, _, _ := config.AutoWire(db)

	if caller != nil {
		return config.SetupRouter(s, h, func(c *gin.Context) {
			c.Set("user", &dto.UserPassport{ID: caller.ID, Username: caller.Username, Role: caller.Role})
		})
	} else {
		return config.SetupRouter(s, h, func(c *gin.Context) {
			c.Set("user", nil)
		})
	}

}
