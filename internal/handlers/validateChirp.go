package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ValidateChirpReq struct {
	Body string `json:"body"`
}

type ChirpCleanedRsp struct {
	CleanedBody string `json:"cleaned_body"`
}

type ChirpErrorRsp struct {
	Error string `json:"error"`
}

var forbiddenWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ValidateChirpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanChirp(req.Body)
	respondWithJSON(w, http.StatusOK, ChirpCleanedRsp{
		CleanedBody: cleanedBody,
	})
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errResp := ChirpErrorRsp{
		Error: msg,
	}

	_ = json.NewEncoder(w).Encode(errResp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	resp, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
