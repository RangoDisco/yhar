package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

func GetOrCreateUser(username, id string) (*models.User, error) {
	existingUser, err := repositories.FindActiveUserByExternalID(id)
	if err == nil {
		return existingUser, err
	}

	model := scrobbleUserToUserModel(username, id)
	err = repositories.PersistUser(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func scrobbleUserToUserModel(username, id string) *models.User {
	return &models.User{
		Username:   username,
		ExternalID: id,
		// TODO: handle enum
		Origin: "SUBSONIC",
	}
}
