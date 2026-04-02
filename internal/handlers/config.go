package handlers

import (
	"sync/atomic"

	"github.com/Raclino/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
}
