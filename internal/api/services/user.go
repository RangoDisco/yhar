package services

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetOrCreateUser(ctx context.Context, username string) (*models.User, error) {
	existingUser, err := s.repo.FindActiveByFilters(ctx, []filters.QueryFilter{
		{Key: "username", Value: username},
	})
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
