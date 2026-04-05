package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

const secret = "super-secret-key"

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	gotUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}

	if gotUserID != userID {
		t.Fatalf("expected userID %v, got %v", userID, gotUserID)
	}
}

func TestValidateJWTExpiredToken(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected expired token to return an error")
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()
	wrongSecret := "wrong-secret-key"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("expected token validated with wrong secret to return an error")
	}
}
