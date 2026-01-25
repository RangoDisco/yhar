package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
)

func PersistScrobble(s *models.Scrobble) error {
	res := database.GetDB().Create(&s)
	return res.Error
}
