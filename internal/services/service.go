package services

import (
	"log/slog"

	"github.com/Raclino/chirpy/internal/database"
)

type Service struct {
	DB        *database.Queries
	Logger    *slog.Logger
	JWTSecret string
	PolkaKey  string
}
