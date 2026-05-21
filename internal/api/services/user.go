package services

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetOrCreateUser(ctx context.Context, rawUsername string) (*models.User, error) {
	var username string
	if rawUsername != "" {
		username = rawUsername
	} else {
		// TODO: Default to first user
		username = "rango"
	}

	existingUser, err := s.repo.FindOneBy(ctx, []repositories.QueryFilter{
		{Key: "username", Value: username},
	}, "Role.Permissions")
	if err == nil {
		return existingUser, err
	}

	model := &models.User{
		Username: username,
		// TODO: handle enum
		Origin: "SUBSONIC",
	}
	err = s.repo.Persist(ctx, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
