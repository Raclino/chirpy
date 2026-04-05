package handlers

import (
	"log/slog"
	"sync/atomic"

	"github.com/Raclino/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits           atomic.Int32
	Db                       *database.Queries
	Platform                 string
	JwtSigningVerifyingToken string
	Logger                   *slog.Logger
}
