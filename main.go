package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/Raclino/chirpy/internal/config"
	"github.com/Raclino/chirpy/internal/handlers"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	apiConfig := config.ApiConfig{
		FileserverHits: atomic.Int32{},
	}
	muxServer := http.NewServeMux()
	muxServer.Handle("/app/", apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	muxServer.HandleFunc("GET /healthz", handlers.HandlerHealth)
	muxServer.HandleFunc("GET /metrics", apiConfig.HandlerMetrics)
	muxServer.HandleFunc("POST /reset", apiConfig.HandlerReset)

	server := &http.Server{
		Handler: muxServer,
		Addr:    ":" + port,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
