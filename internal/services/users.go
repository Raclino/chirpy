package services

import (
	"context"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
	"github.com/Raclino/chirpy/internal/database"
	"github.com/google/uuid"
)

var oneHourExpiresInSec = 3600

func (s *Service) CreateUser(ctx context.Context, email, password string) (database.CreateUserRow, error) {
	now := time.Now().UTC()

	hashedPwd, err := auth.HashPassword(password)
	if err != nil {
		return database.CreateUserRow{}, err
	}

	params := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      now,
		UpdatedAt:      now,
		Email:          email,
		HashedPassword: hashedPwd,
	}

	user, err := s.DB.CreateUser(ctx, params)
	if err != nil {
		return database.CreateUserRow{}, err
	}

	s.Logger.Info("user created", "user_id", user.ID.String(), "email", user.Email)

	return user, nil
}
