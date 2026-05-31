package factories

import (
	"context"
	"testing"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func SeedUser(t *testing.T, db *gorm.DB, name, password, role string, isPublic bool) models.User {
	t.Helper()

	var roleId int64
	switch role {
	case "ADMIN":
		roleId = 1
	default:
		roleId = 2
	}

	ctx := context.Background()
	rawUser := models.User{
		Username: name,
		Password: password,
		RoleID:   roleId,
		IsPublic: isPublic,
	}

	err := db.WithContext(ctx).Create(&rawUser).Error
	require.NoError(t, err)

	var u models.User
	err = db.WithContext(ctx).Preload("Role.Permissions").Where("id = ?", rawUser.ID).Find(&u).Error
	require.NoError(t, err)

	return u
}
