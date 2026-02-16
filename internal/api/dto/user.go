package dto

import "github.com/rangodisco/yhar/internal/api/models"

type UserPassport struct {
	ID       int64       `json:"id"`
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
}
