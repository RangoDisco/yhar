package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/filters"
)

func FindActiveUserByFilters(filters []filters.QueryFilter) (*models.User, error) {
	var u models.User
	query := database.GetDB().Preload("Role.Permissions")

	for _, filter := range filters {
		query.Where(filter.Key+" = ?", filter.Value)
	}

	err := query.First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func PersistUser(user *models.User) error {
	res := database.GetDB().Create(user)
	return res.Error
}
