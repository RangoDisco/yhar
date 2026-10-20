package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/api/models"
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

	err = addSearchTrigger(db)
	if err != nil {
		return nil, err
	}

	err = addViews(db)
	if err != nil {
		return nil, err
	}

	return db, err
}

func addSearchTrigger(db *gorm.DB) error {
	err := db.Exec(`
				CREATE OR REPLACE FUNCTION artist_search_trigger() RETURNS trigger AS
				$$
				BEGIN
					NEW.search_vector := to_tsvector('english', COALESCE(NEW.name, ''));
					RETURN NEW;
				END;
				$$ LANGUAGE plpgsql;
				
				CREATE OR REPLACE TRIGGER trg_artist_search
					BEFORE INSERT OR UPDATE
					ON artists
					FOR EACH ROW
				EXECUTE FUNCTION artist_search_trigger();
				
				CREATE OR REPLACE FUNCTION content_search_trigger() RETURNS trigger AS
				$$
				BEGIN
					NEW.search_vector := to_tsvector('english', COALESCE(NEW.title, ''));
					RETURN NEW;
				END;
				$$ LANGUAGE plpgsql;
				
				CREATE OR REPLACE TRIGGER trg_artist_search
					BEFORE INSERT OR UPDATE
					ON artists
					FOR EACH ROW
				EXECUTE FUNCTION artist_search_trigger();
				
				CREATE OR REPLACE TRIGGER trg_album_search
					BEFORE INSERT OR UPDATE
					ON albums
					FOR EACH ROW
				EXECUTE FUNCTION content_search_trigger();
				
				CREATE OR REPLACE TRIGGER trg_track_search
					BEFORE INSERT OR UPDATE
					ON tracks
					FOR EACH ROW
				EXECUTE FUNCTION content_search_trigger();
`).Error
	if err != nil {
		return err
	}

	return nil
}

func addViews(db *gorm.DB) error {
	searchQuery := db.Raw(`
					SELECT 'artists' as content_type, a.id as content_id, search_vector as value, a.name as title, i.path, i.domain
					FROM artists a
							 LEFT JOIN images i ON i.id = a.picture_id
					WHERE a.deleted_at IS null
					UNION
					SELECT 'albums' AS content_type, al.id as content_id, search_vector as value, al.title, i.path, i.domain
					FROM albums al
							 LEFT JOIN images i ON i.id = al.picture_id
					WHERE al.deleted_at IS null
					UNION
					SELECT 'tracks' AS content_type, tr.id as content_id, tr.search_vector as value, tr.title, i.path, i.domain
					FROM tracks tr
							 LEFT JOIN albums al ON al.id = tr.album_id
							 LEFT JOIN images i ON i.id = al.picture_id
					WHERE tr.deleted_at IS null
`)

	err := db.Migrator().CreateView("content_search", gorm.ViewOption{Query: searchQuery, Replace: true})
	if err != nil {
		return err
	}
	return nil
}
