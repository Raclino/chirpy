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
	if err := run(ctx, os.Getenv, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string, stdout io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		return err
	}

	platform := getenv("PLATFORM")
	jwtSigningVerifyingToken := getenv("JWT_SIGNING_VERIFYING")
	dbURL := getenv("DB_URL")
	polkaKey := getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	dbQueries := database.New(db)
	logger := slog.New(slog.NewTextHandler(stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	apiConfig := &handlers.ApiConfig{
		FileserverHits:           atomic.Int32{},
		Db:                       dbQueries,
		Platform:                 platform,
		JwtSigningVerifyingToken: jwtSigningVerifyingToken,
		Logger:                   logger,
		PolkaKey:                 polkaKey,
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
