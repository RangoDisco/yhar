package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
)

type UserService struct {
	uRepo *repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{uRepo: repository}
}

func (s *UserService) GetOrCreateUser(username string) (*models.User, error) {
	uFilters := []filters.QueryFilter{
		{Key: "username", Value: username},
	}
	existingUser, err := s.uRepo.FindActiveUserByFilters(uFilters)
	if err == nil {
		return existingUser, err
	}

	model := s.scrobbleUserToUserModel(username)
	err = s.uRepo.PersistUser(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *UserService) scrobbleUserToUserModel(username string) *models.User {
	return &models.User{
		Username: username,
		// TODO: handle enum
		Origin: "SUBSONIC",
	}
}
