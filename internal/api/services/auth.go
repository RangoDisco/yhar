package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/auth"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(u *repositories.UserRepository) *AuthService {
	return &AuthService{repo: u}
}

func (s *AuthService) EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ComparePassword checks if the given plain password correspond to the given hash
func (s *AuthService) ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateToken creates a JWT with user's name as claim
func (s *AuthService) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("unable to creat token: %w", err)
	}

	return tokenString, nil
}

// HandleUserLogin receives an auth.LoginRequest
// tries to find user by its username
// compares the passwords and creates a token
func (s *AuthService) HandleUserLogin(ctx context.Context, request auth.LoginRequest) (string, error) {

	user, err := s.repo.FindActiveByFilters(ctx, []filters.QueryFilter{
		{Key: "username", Value: request.Username},
	})
	if err != nil {
		return "", err
	}

	success := s.ComparePassword(request.Password, user.Password)
	if !success {
		return "", errors.New("invalid password")
	}

	token, err := s.CreateToken(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ParseToken parses the given tokenString with
func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	return token, nil
}

// GetUserFromToken uses the username in the claims to find a user by its username
func (s *AuthService) GetUserFromToken(ctx context.Context, token *jwt.Token) (*models.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token is invalid or expired")
	}

	rawUsername, ok := claims["username"]
	if !ok {
		return nil, errors.New("token missing username claim")
	}
	username, ok := rawUsername.(string)
	if !ok || username == "" {
		return nil, errors.New("username claim is not a valid string")
	}

	user, err := s.repo.FindActiveByFilters(ctx, []filters.QueryFilter{
		{Key: "username", Value: username},
	})

	if err != nil {
		return nil, fmt.Errorf("unable to get user from token: %w", err)
	}

	return user, nil
}
