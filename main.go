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
	const filePathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

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
	}
	muxServer := http.NewServeMux()
	muxServer.Handle("/app/", apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	muxServer.HandleFunc("GET /api/healthz", handlers.HandlerHealth)
	muxServer.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidateChirp)
	muxServer.HandleFunc("POST /admin/reset", apiConfig.HandlerReset)
	muxServer.HandleFunc("GET /admin/metrics", apiConfig.HandlerMetrics)

	server := &http.Server{
		Handler: muxServer,
		Addr:    ":" + port,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
