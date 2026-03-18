package database

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var logLevel logger.LogLevel

func SetupDatabase() (*gorm.DB, error) {
	var ginMode = os.Getenv("GIN_MODE")
	switch ginMode {
	case gin.DebugMode:
		logLevel = logger.Info
		return InitDatabase()
	case gin.TestMode:
		logLevel = logger.Info
		panic("not implemented yet")
	case gin.ReleaseMode:
	default:
		logLevel = logger.Silent
		return InitDatabase()
	}
	return nil, nil
}

func InitDatabase() (*gorm.DB, error) {
	name := os.Getenv("YHAR_DB_NAME")
	user := os.Getenv("YHAR_DB_USER")
	password := os.Getenv("YHAR_DB_PASSWORD")
	host := os.Getenv("YHAR_DB_HOST")
	port := os.Getenv("YHAR_DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
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

	return db, err
}
