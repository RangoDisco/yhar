package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
)

func GetOrCreateUser(username string) (*models.User, error) {
	uFilters := []filters.QueryFilter{
		{Key: "username", Value: username},
	}
	existingUser, err := repositories.FindActiveUserByFilters(uFilters)
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
