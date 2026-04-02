package handlers

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits = atomic.Int32{}

	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
	}

	if err := cfg.Db.DeleteAllUsers(r.Context()); err != nil {
		fmt.Println("couldn't delete all users: %w", err)
	}
}
