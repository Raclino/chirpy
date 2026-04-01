package handlers

import (
	"net/http"
	"sync/atomic"

	"github.com/Raclino/chirpy/internal/config"
)

func (cfg *config.ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits = atomic.Int32{}
}
