package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filePathRoot = "."
	const port = "8080"
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	muxServer := http.NewServeMux()
	muxServer.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	muxServer.HandleFunc("/healthz", handlerHealth)
	muxServer.HandleFunc("/metrics", apiConfig.handlerMetrics)
	muxServer.HandleFunc("/reset", apiConfig.handlerReset)

	server := &http.Server{
		Handler: muxServer,
		Addr:    ":" + port,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf("Hits: %d\n", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
