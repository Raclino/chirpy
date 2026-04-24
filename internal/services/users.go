package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
	"github.com/Raclino/chirpy/internal/database"
	"github.com/google/uuid"
)

var oneHourExpiresInSec = 3600

type LoginResult struct {
	User         database.User
	AccessToken  string
	RefreshToken string
}

type UpdateUserResult struct {
	User        database.UpdateUserPwdEmailRow
	AccessToken string
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (database.CreateUserRow, error) {
	now := time.Now().UTC()

	hashedPwd, err := auth.HashPassword(password)
	if err != nil {
		return database.CreateUserRow{}, fmt.Errorf("hash password: %w", err)
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

func (s *Service) UserLogin(ctx context.Context, email, password string) (LoginResult, error) {
	user, err := s.DB.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	ok, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		return LoginResult{}, fmt.Errorf("check password hash: %w", err)
	}
	if !ok {
		return LoginResult{}, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	refreshToken := auth.MakeRefreshToken()

	createdRefreshToken, err := s.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.AddDate(0, 0, 60),
		RevokedAt: sql.NullTime{},
		UserID:    user.ID,
	})
	if err != nil {
		return LoginResult{}, fmt.Errorf("create refresh token: %w", err)
	}

	accessToken, err := auth.MakeJWT(user.ID, s.JWTSecret, time.Hour)
	if err != nil {
		return LoginResult{}, fmt.Errorf("make jwt: %w", err)
	}

	return LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: createdRefreshToken,
	}, nil
}

func (s *Service) UpdateUsers(ctx context.Context, jwtToken, email, password string) (UpdateUserResult, error) {
	userID, err := auth.ValidateJWT(jwtToken, s.JWTSecret)
	if err != nil {
		return UpdateUserResult{}, ErrUnauthorized
	}

	hashedPwd, err := auth.HashPassword(password)
	if err != nil {
		return UpdateUserResult{}, fmt.Errorf("hash password: %w", err)
	}

	params := database.UpdateUserPwdEmailParams{
		ID:             userID,
		Email:          email,
		HashedPassword: hashedPwd,
		UpdatedAt:      time.Now().UTC(),
	}

	updatedUser, err := s.DB.UpdateUserPwdEmail(ctx, params)
	if err != nil {
		return UpdateUserResult{}, fmt.Errorf("update user: %w", err)
	}

	return UpdateUserResult{
		User:        updatedUser,
		AccessToken: jwtToken,
	}, nil
}
