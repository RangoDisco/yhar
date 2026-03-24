package testutils

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/config"
	"github.com/rangodisco/yhar/internal/api/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	db, _ := gorm.Open(postgres.Open(dsn))

	tx := db.Begin()
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
