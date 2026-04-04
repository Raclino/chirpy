package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorRsp struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errResp := ErrorRsp{
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
	resp = append(resp, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
