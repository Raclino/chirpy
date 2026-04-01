package handlers

import (
	"net/http"
	"sync/atomic"
)

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
}
