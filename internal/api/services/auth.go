package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/auth"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ComparePassword ensures if the given plain password correspond to the given hash
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
	user, err := repositories.FindActiveByUsername(request.Username)
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
