package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ValidateChirpReq struct {
	Body string `json:"body"`
}

type ChirpValidRsp struct {
	Valid bool `json:"valid"`
}

type ChirpErrorRsp struct {
	Error string `json:"error"`
}

func HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	req := ValidateChirpReq{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")

	}
	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	} else {
		respondWithJSON(w, http.StatusOK)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errResp := ChirpErrorRsp{
		Error: msg,
	}
	resp, err := json.Marshal(errResp)
	if err != nil {
		fmt.Println("Couldn't marshal the msg")
	}
	w.Write(resp)
}

func respondWithJSON(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	valid := ChirpValidRsp{
		Valid: true,
	}
	resp, err := json.Marshal(valid)
	if err != nil {
		fmt.Println("Couldn't marshal the msg")
	}
	w.Write(resp)
}
