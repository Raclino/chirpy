package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
	"github.com/Raclino/chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	requestRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		cfg.Logger.Warn("missing or invalid refresh bearer token",
			"path", r.URL.Path,
			"method", r.Method,
			"error", err,
		)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := cfg.Db.GetTokenByRefreshToken(r.Context(), requestRefreshToken)
	if err != nil {
		cfg.Logger.Warn("refresh token not found",
			"path", r.URL.Path,
			"method", r.Method,
			"error", err,
		)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if token.ExpiresAt.Before(time.Now().UTC()) {
		cfg.Logger.Warn("refresh token expired",
			"path", r.URL.Path,
			"method", r.Method,
			"user_id", token.UserID.String(),
			"expires_at", token.ExpiresAt,
		)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if token.RevokedAt.Valid {
		cfg.Logger.Warn("refresh token revoked",
			"path", r.URL.Path,
			"method", r.Method,
			"user_id", token.UserID.String(),
			"revoked_at", token.RevokedAt.Time,
		)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	newToken, err := auth.MakeJWT(
		token.UserID,
		cfg.JwtSigningVerifyingToken,
		time.Duration(oneHourExpiresInSec)*time.Second,
	)
	if err != nil {
		cfg.Logger.Error("failed to create access token from refresh token",
			"path", r.URL.Path,
			"method", r.Method,
			"user_id", token.UserID.String(),
			"error", err,
		)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}

	cfg.Logger.Info("access token refreshed",
		"path", r.URL.Path,
		"method", r.Method,
		"user_id", token.UserID.String(),
	)

	respondWithJSON(w, http.StatusOK, RefreshToken{Token: newToken})
}

func (cfg *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	requestRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		cfg.Logger.Warn("missing or invalid refresh bearer token for revoke",
			"path", r.URL.Path,
			"method", r.Method,
			"error", err,
		)
		respondWithError(w, http.StatusUnauthorized, "Couldn't get refresh_token")
		return
	}

	now := time.Now().UTC()
	params := database.RevokeRefreshTokenParams{
		Token:     requestRefreshToken,
		UpdatedAt: now,
		RevokedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	revokedToken, err := cfg.Db.RevokeRefreshToken(r.Context(), params)
	if err != nil {
		cfg.Logger.Warn("failed to revoke refresh token",
			"path", r.URL.Path,
			"method", r.Method,
			"error", err,
		)
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke refresh_token")
		return
	}

	cfg.Logger.Info("refresh token revoked",
		"path", r.URL.Path,
		"method", r.Method,
		"user_id", revokedToken.UserID.String(),
		"revoked_at", revokedToken.RevokedAt.Time,
	)

	w.WriteHeader(http.StatusNoContent)
}
