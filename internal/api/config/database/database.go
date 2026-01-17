package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func SetupDatabase() error {
	var ginMode = os.Getenv("GIN_MODE")
	switch ginMode {
	case "release":
	default:
		err = InitDatabase()
	}

	return err
}

func InitDatabase() error {
	name := os.Getenv("YHAR_DB_NAME")
	user := os.Getenv("YHAR_DB_USER")
	password := os.Getenv("YHAR_DB_PASSWORD")
	host := os.Getenv("YHAR_DB_HOST")
	port := os.Getenv("YHAR_DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return err
}

func GetDB() *gorm.DB {
	return db
}
