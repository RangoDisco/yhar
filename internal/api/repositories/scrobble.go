package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/metadata/config/database"
)

func PersistScrobble(s *models.Scrobble) error {
	res := database.GetDB().Create(&s)
	return res.Error
}
