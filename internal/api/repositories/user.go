package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
)

func FindActiveByUsername(username string) (*models.User, error) {
	var u models.User
	err := database.GetDB().Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func PersistUser(user *models.User) error {
	res := database.GetDB().Create(user)
	return res.Error
}
