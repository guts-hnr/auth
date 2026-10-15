package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	return err == nil
}

func GenerateToken(userId int, secret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userId),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenstring string, secret string) (int, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenstring,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
