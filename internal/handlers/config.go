package handlers

import (
	"log/slog"
	"sync/atomic"

	"github.com/Raclino/chirpy/internal/database"
	"github.com/Raclino/chirpy/internal/services"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	//TODO supprimé une fois tous les handlers migré via service
	Db       *database.Queries
	Platform string
	Logger   *slog.Logger

	JwtSigningVerifyingToken string
	PolkaKey                 string

	Service *services.Service
}
