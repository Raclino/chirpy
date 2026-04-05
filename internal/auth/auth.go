package auth

import (
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashedPwd, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("couldn't hash the pwd: %w", err)
	}
	return hashedPwd, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, _, err := argon2id.CheckHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("couldn't check pwd with hash: %w", err)
	}

	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", fmt.Errorf("error while signing the token: %w", err)
	}

	fmt.Println(ss, err)
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.Claims{}
	keyFunc := jwt.Keyfunc{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return uuid.UUIDs, nil
	}

	return token.Raw, nil
}
