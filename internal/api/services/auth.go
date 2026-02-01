package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/auth"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ComparePassword checks if the given plain password correspond to the given hash
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateToken creates a JWT with user's name as claim
func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// HandleUserLogin receives an auth.LoginRequest
// tries to find user by its username
// compares the passwords and creates a token
func HandleUserLogin(request auth.LoginRequest) (string, error) {
	uFilters := []filters.QueryFilter{
		{Key: "username", Value: request.Username},
	}
	user, err := repositories.FindActiveUserByFilters(uFilters)
	if err != nil {
		return "", err
	}

	success := ComparePassword(request.Password, user.Password)
	if !success {
		return "", errors.New("invalid password")
	}

	token, err := CreateToken(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	return token, err
}

// GetUserFromToken uses the username in the claims to find a user by its username
func GetUserFromToken(token *jwt.Token) (*models.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	username := claims["username"].(string)

	uFilters := []filters.QueryFilter{
		{Key: "username", Value: username},
	}
	user, err := repositories.FindActiveUserByFilters(uFilters)

	return user, err
}
