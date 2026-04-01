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
		fileserverHits: atomic.Int32{},
	}
	muxServer := http.NewServeMux()
	muxServer.Handle("/app/", ApiConfig.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	muxServer.HandleFunc("GET /healthz", handlers.HandlerHealth)
	muxServer.HandleFunc("GET /metrics", ApiConfig.HandlerMetrics)
	muxServer.HandleFunc("POST /reset", ApiConfig.HandlerReset)

	server := &http.Server{
		Handler: muxServer,
		Addr:    ":" + port,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
