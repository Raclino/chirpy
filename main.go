package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Raclino/chirpy/internal/database"
	"github.com/Raclino/chirpy/internal/handlers"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	const filePathRoot = "."
	const port = "8080"

	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	apiConfig := &handlers.ApiConfig{
		FileserverHits: atomic.Int32{},
		Db:             dbQueries,
		Platform:       platform,
	}

	muxServer := http.NewServeMux()
	muxServer.Handle("/app/", apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	muxServer.HandleFunc("GET /api/healthz", handlers.HandlerGetHealth)
	muxServer.HandleFunc("POST /api/users", apiConfig.HandlerCreateUsers)
	muxServer.HandleFunc("POST /api/login", apiConfig.HandlerUserLogin)
	muxServer.HandleFunc("POST /api/chirps", apiConfig.HandlerCreateChirps)
	muxServer.HandleFunc("GET /api/chirps", apiConfig.HandlerGetChirps)
	muxServer.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.HandlerGetChirpsByID)
	muxServer.HandleFunc("POST /admin/reset", apiConfig.HandlerReset)
	muxServer.HandleFunc("GET /admin/metrics", apiConfig.HandlerGetMetrics)

	server := &http.Server{
		Handler: muxServer,
		Addr:    ":" + port,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
