package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
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

	if !match {
		return false, fmt.Errorf("pwd is incorrect")
	}

	return match, nil
}
