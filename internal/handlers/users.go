package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
	"github.com/Raclino/chirpy/internal/database"
	"github.com/Raclino/chirpy/internal/services"
	"github.com/google/uuid"
)

type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type RefreshToken struct {
	Token string `json:"token"`
}

type WebHooksReq struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

var oneHourExpiresInSec = 3600

func (cfg *ApiConfig) HandleCreateUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Warn("invalid create user request body", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdUser, err := cfg.Service.CreateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		cfg.Logger.Warn("failed to create user", "path", r.URL.Path, "method", r.Method, "email", req.Email, "error", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't create user")
		return
	}

	user := User{
		ID:          createdUser.ID,
		CreatedAt:   createdUser.CreatedAt,
		UpdatedAt:   createdUser.UpdatedAt,
		Email:       createdUser.Email,
		IsChirpyRed: createdUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *ApiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Warn("invalid login request body", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := cfg.Service.UserLogin(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			cfg.Logger.Warn("login failed", "path", r.URL.Path, "method", r.Method, "email", req.Email)
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		cfg.Logger.Error("login failed", "path", r.URL.Path, "method", r.Method, "email", req.Email, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't login user")
		return
	}

	cfg.Logger.Info("user logged in", "user_id", result.User.ID.String(), "email", result.User.Email)

	userResponse := User{
		ID:           result.User.ID,
		CreatedAt:    result.User.CreatedAt,
		UpdatedAt:    result.User.UpdatedAt,
		Email:        result.User.Email,
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
		IsChirpyRed:  result.User.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, userResponse)
}

func (cfg *ApiConfig) HandleUpdateUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		cfg.Logger.Warn("missing or invalid Authorization header", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.JwtSigningVerifyingToken)
	if err != nil {
		cfg.Logger.Warn("invalid jwt token", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Warn("invalid update user request body", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hashedPwd, err := auth.HashPassword(req.Password)
	if err != nil {
		cfg.Logger.Error("failed to hash password", "path", r.URL.Path, "method", r.Method, "email", req.Email, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	params := database.UpdateUserPwdEmailParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: hashedPwd,
		UpdatedAt:      time.Now().UTC(),
	}

	updatedUser, err := cfg.Db.UpdateUserPwdEmail(r.Context(), params)
	if err != nil {
		cfg.Logger.Error("failed to update user", "path", r.URL.Path, "method", r.Method, "user_id", userID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
	}

	user := User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
		Token:     jwtToken,
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *ApiConfig) HandlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		cfg.Logger.Error("invalid apiKey webhooks user ", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid request header")
		return
	}

	if apiKey != cfg.PolkaKey {
		cfg.Logger.Warn("invalid webhooks user request body", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid request header")
	}

	var req WebHooksReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Warn("invalid webhooks user request body", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	dbUser, err := cfg.Db.UpdateUserChirpyRedMembership(r.Context(), database.UpdateUserChirpyRedMembershipParams{
		ID:        req.Data.UserID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	fmt.Printf("dbUser: %v+", dbUser)

	w.WriteHeader(http.StatusNoContent)
}
