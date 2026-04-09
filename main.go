package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/Raclino/chirpy/internal/database"
	"github.com/Raclino/chirpy/internal/handlers"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

const (
	filePathRoot = "."
	port         = "8080"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		return err
	}

	platform := os.Getenv("PLATFORM")
	jwtSigningVerifyingToken := os.Getenv("JWT_SIGNING_VERIFYING")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	dbQueries := database.New(db)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	apiConfig := &handlers.ApiConfig{
		FileserverHits:           atomic.Int32{},
		Db:                       dbQueries,
		Platform:                 platform,
		JwtSigningVerifyingToken: jwtSigningVerifyingToken,
		Logger:                   logger,
	}

	srv := NewServer(apiConfig)
	httpServer := &http.Server{
		Handler: srv,
		Addr:    ":" + port,
	}

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	logger.Info("starting server", "addr", httpServer.Addr)

	return err
}
