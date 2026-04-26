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

// CreateTokens creates a JWT and its refresh with user's name as claim
func (s *AuthService) CreateTokens(username, role string) (string, string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Minute * 60).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", "", fmt.Errorf("unable to create token: %w", err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if err != nil {
		return "", "", fmt.Errorf("unable to create refresh token: %w", err)
	}

	return tokenString, refreshTokenString, nil
}

func (s *AuthService) RefreshToken(t string) (string, error) {
	secret := os.Getenv("REFRESH_SECRET")
	token, err := ParseToken(t, secret)
	if err != nil {
		return "", fmt.Errorf("unable to parse token: %w", err)
	}

	newToken, _, err := s.CreateTokens(token.Claims.(jwt.MapClaims)["username"].(string), token.Claims.(jwt.MapClaims)["role"].(string))
	if err != nil {
		return "", fmt.Errorf("unable to generate new token from refresh: %w", err)
	}

	return newToken, nil
}

// HandleUserLogin receives an auth.LoginRequest
// tries to find user by its username
// compares the passwords and creates a token
func (s *AuthService) HandleUserLogin(ctx context.Context, request auth.LoginRequest) (string, string, error) {

	user, err := s.repo.FindActiveByFilters(ctx, []filters.QueryFilter{
		{Key: "username", Value: request.Username},
	})
	if err != nil {
		return "", "", err
	}

	success := s.ComparePassword(request.Password, user.Password)
	if !success {
		return "", "", errors.New("invalid password")
	}

	token, refresh, err := s.CreateTokens(user.Username, user.Role.Name)
	if err != nil {
		return "", "", err
	}

	return token, refresh, nil
}

// ParseToken parses the given tokenString with
func ParseToken(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
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
