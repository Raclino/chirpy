package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
	"github.com/Raclino/chirpy/internal/database"
	"github.com/google/uuid"
)

type CreateChirpReq struct {
	Body string `json:"body"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

var forbiddenWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func (cfg *ApiConfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")

	if authorIDStr != "" {
		authorID, err := uuid.Parse(authorIDStr)
		if err != nil {
			cfg.Logger.Warn("invalid author_id query param",
				"path", r.URL.Path,
				"method", r.Method,
				"author_id", authorIDStr,
				"error", err,
			)
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}

		dbChirps, err := cfg.Db.GetChirpsByAuthorID(r.Context(), authorID)
		if err != nil {
			cfg.Logger.Error("failed to get chirps by author id",
				"path", r.URL.Path,
				"method", r.Method,
				"author_id", authorID.String(),
				"error", err,
			)
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
			return
		}

		chirps := make([]Chirp, 0, len(dbChirps))
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, mapDBChirpToResponse(dbChirp))
		}

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	dbChirps, err := cfg.Db.GetChirps(r.Context())
	if err != nil {
		cfg.Logger.Error("failed to get chirps", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, mapDBChirpToResponse(dbChirp))
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *ApiConfig) HandleGetChirpsByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		cfg.Logger.Warn("invalid chirp id in path", "path", r.URL.Path, "method", r.Method, "chirp_id", chirpIDStr, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}

	dbChirp, err := cfg.Db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		cfg.Logger.Warn("chirp not found", "path", r.URL.Path, "method", r.Method, "chirp_id", chirpID.String(), "error", err)
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	respondWithJSON(w, http.StatusOK, mapDBChirpToResponse(dbChirp))
}

func (cfg *ApiConfig) HandleCreateChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateChirpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Warn("invalid create chirp request body", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		cfg.Logger.Warn("missing or invalid bearer token", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.JwtSigningVerifyingToken)
	if err != nil {
		cfg.Logger.Warn("invalid jwt token", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if len(req.Body) > 140 {
		cfg.Logger.Warn("chirp too long", "path", r.URL.Path, "method", r.Method, "user_id", userID.String(), "body_length", len(req.Body))
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanChirp(req.Body)
	now := time.Now()

	params := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Body:      cleanedBody,
		UserID:    userID,
	}

	createdChirp, err := cfg.Db.CreateChirp(r.Context(), params)
	if err != nil {
		cfg.Logger.Error("failed to create chirp", "path", r.URL.Path, "method", r.Method, "user_id", userID.String(), "error", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	cfg.Logger.Info("chirp created", "chirp_id", createdChirp.ID.String(), "user_id", createdChirp.UserID.String())

	respondWithJSON(w, http.StatusCreated, mapDBChirpToResponse(createdChirp))
}
func (cfg *ApiConfig) HandleDeleteChirpsByID(w http.ResponseWriter, r *http.Request) {
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

	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		cfg.Logger.Warn("invalid chirp queryParams", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	dbChirp, err := cfg.Db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		cfg.Logger.Warn("error getting chirp by id", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusNotFound, "")
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	if err := cfg.Db.DeleteChirpByID(r.Context(), dbChirp.ID); err != nil {
		cfg.Logger.Warn("error deleting chirp by id", "path", r.URL.Path, "method", r.Method, "error", err)
		respondWithError(w, http.StatusForbidden, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func cleanChirp(body string) string {
	words := strings.Split(body, " ")

	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, found := forbiddenWords[lowerWord]; found {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func mapDBChirpToResponse(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}
