package database

import (
	"fmt"
	"os"
	"time"

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
	default:
		logLevel = logger.Silent
		return InitDatabase()
	}
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
		Logger:         logger.Default.LogMode(logLevel),
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(20)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

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
