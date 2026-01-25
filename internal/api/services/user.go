package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

func GetOrCreateUser(username string) (*models.User, error) {
	existingUser, err := repositories.FindActiveByUsername(username)
	if err == nil {
		return existingUser, err
	}

	model := scrobbleUserToUserModel(username)
	err = repositories.PersistUser(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func scrobbleUserToUserModel(username string) *models.User {
	return &models.User{
		Username: username,
		// TODO: handle enum
		Origin: "SUBSONIC",
	}
}
