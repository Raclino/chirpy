package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
	"github.com/Raclino/chirpy/internal/database"
	"github.com/google/uuid"
)

type UserReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ExpiresInSec int    `json:"expires_in_seconds"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token,omitempty"`
}

var defaultExpiresInSec = 3600

func (cfg *ApiConfig) HandlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Warn("invalid create user request body", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	now := time.Now()
	hashedPwd, err := auth.HashPassword(req.Password)
	if err != nil {
		cfg.Logger.Error("failed to hash password", "path", r.URL.Path, "method", r.Method, "email", req.Email, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	newUser := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      now,
		UpdatedAt:      now,
		Email:          req.Email,
		HashedPassword: hashedPwd,
	}

	createdUser, err := cfg.Db.CreateUser(r.Context(), newUser)
	if err != nil {
		cfg.Logger.Warn("failed to create user", "path", r.URL.Path, "method", r.Method, "email", req.Email, "error", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't create user")
		return
	}

	cfg.Logger.Info("user created", "user_id", createdUser.ID.String(), "email", createdUser.Email)

	user := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *ApiConfig) HandlerUserLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Warn("invalid login request body", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := cfg.Db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		cfg.Logger.Warn("login failed: user lookup", "email", req.Email, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	isPwdValid, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		cfg.Logger.Error("failed to verify password hash", "user_id", user.ID.String(), "email", req.Email, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if !isPwdValid {
		cfg.Logger.Warn("login failed: invalid password", "user_id", user.ID.String(), "email", req.Email)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expiresInSec := defaultExpiresInSec
	if req.ExpiresInSec > 0 && req.ExpiresInSec < defaultExpiresInSec {
		expiresInSec = req.ExpiresInSec
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.JwtSigningVerifyingToken, time.Duration(expiresInSec)*time.Second)
	if err != nil {
		cfg.Logger.Error("failed to create jwt", "user_id", user.ID.String(), "email", user.Email, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}

	cfg.Logger.Info("user logged in", "user_id", user.ID.String(), "email", user.Email)

	userResponse := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwtToken,
	}

	respondWithJSON(w, http.StatusOK, userResponse)
}
