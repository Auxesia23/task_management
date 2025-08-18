package auth

import (
	"errors"
	"os"
	"time"

	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates an access token for a user
func GenerateAccessToken(user *models.User) (string, error) {
	var secret []byte
	secret = []byte(os.Getenv("JWT_ACCESS_SECRET"))

	claims := &dto.AccessTokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			Issuer:    "Auxesia",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", errors.New("error signing token")
	}
	return tokenString, nil
}

func ValidateAccessToken(tokenString string) (*dto.AccessTokenClaims, error) {
	var secret []byte

	secret = []byte(os.Getenv("JWT_ACCESS_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &dto.AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, errors.New("error parsing token")
	}

	if claims, ok := token.Claims.(*dto.AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GenerateRefreshToken(user *models.User) (string, error) {
	var secret []byte
	secret = []byte(os.Getenv("JWT_REFRESH_SECRET"))

	claims := &dto.RefreshTokenClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			Issuer:    "Auxesia",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", errors.New("error signing token")
	}
	return tokenString, nil
}

func ValidateRefreshToken(tokenString string) (*dto.RefreshTokenClaims, error) {
	var secret []byte

	secret = []byte(os.Getenv("JWT_REFRESH_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &dto.RefreshTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, errors.New("error parsing token")
	}

	if claims, ok := token.Claims.(*dto.RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
