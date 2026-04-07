package handlers

import (
	"net/http"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	requestRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get refresh_token")
		return
	}

	token, err := cfg.Db.GetTokenByRefreshToken(r.Context(), requestRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh_token doesn't exist in DB")
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		cfg.Logger.Error("refresh_token", "user_id", token.UserID, "error", err)
		respondWithError(w, http.StatusUnauthorized, "refresh_token is expired")
		return
	}

	if token.RevokedAt.Valid {
		cfg.Logger.Error("refresh_token", "user_id", requestRefreshToken, "error", err)
		respondWithError(w, http.StatusUnauthorized, "refresh_token was revoked")
		return
	}

	newToken, err := auth.MakeJWT(token.UserID, cfg.JwtSigningVerifyingToken, time.Duration(oneHourExpiresInSec)*time.Second)
	if err != nil {
		cfg.Logger.Error("failed to create jwt", "user_id", token.UserID.String(), "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshToken{Token: newToken})
}

func (cfg *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {

	return
}
