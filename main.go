package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/Raclino/chirpy/internal/handlers"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	apiConfig := handlers.ApiConfig{
		FileserverHits: atomic.Int32{},
	}
	muxServer := http.NewServeMux()
	muxServer.Handle("/app/", apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	muxServer.HandleFunc("GET /api/healthz", handlers.HandlerHealth)
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
