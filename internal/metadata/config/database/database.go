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
	var err error
	var ginMode = os.Getenv("APP_ENV")
	switch ginMode {
	default:
		err = InitDatabase()
	}

	return err
}

func InitDatabase() error {
	// Open a database connection
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return err
}

func GetDB() *gorm.DB {
	return db
}
