package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Raclino/chirpy/internal/auth"
	"github.com/Raclino/chirpy/internal/database"
	"github.com/google/uuid"
)

type ValidateChirpReq struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type ChirpCleanedRsp struct {
	CleanedBody string `json:"cleaned_body"`
}

var forbiddenWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func (cfg *ApiConfig) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.Db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
func (cfg *ApiConfig) HandlerGetChirpsByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	println(chirpID)

	dbChirp, err := cfg.Db.GetChirpByID(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *ApiConfig) HandlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ValidateChirpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("jwtToken Error: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDFromJwtToken, err := auth.ValidateJWT(jwtToken, cfg.JwtSigningVerifyingToken)
	if err != nil {
		fmt.Printf("ValidateJWT: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, token not valid")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanChirp(req.Body)

	now := time.Now()
	chirpParams := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Body:      cleanedBody,
		UserID:    userIDFromJwtToken,
	}

	createdChirp, err := cfg.Db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		fmt.Printf("couldn't create chirp in db: %v\n", err)

		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	chirp := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirp)

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
