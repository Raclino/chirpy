package auth

import (
	"fmt"
	"net/http"
	"strings"
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
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("couldn't sign token: %w", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("couldn't parse token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subject uuid: %w", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	fmt.Printf("headers: %v\n", headers)
	tokenString := headers.Get("TOKEN_STRING")
	if tokenString == "" {
		return "", fmt.Errorf("couldn't get TOKEN_STRING from header: %s", tokenString)

	}
	token, found := strings.CutPrefix(tokenString, "Bearer")
	if !found {
		return "", fmt.Errorf("couldn't cut prefix Bearer")
	}
	return strings.TrimSpace(token), nil
}
