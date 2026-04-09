package main

import (
	"net/http"

	"github.com/Raclino/chirpy/internal/handlers"
)

func addRoutes(
	mux *http.ServeMux,
	apiConfig *handlers.ApiConfig,
) {

	mux.Handle("/app/", apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlers.HandlerGetHealth)
	mux.HandleFunc("POST /api/users", apiConfig.HandlerCreateUsers)
	mux.HandleFunc("PUT /api/users", apiConfig.HandlerUpdateUsers)
	mux.HandleFunc("POST /api/login", apiConfig.HandlerUserLogin)
	mux.HandleFunc("POST /api/chirps", apiConfig.HandlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiConfig.HandlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.HandlerGetChirpsByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.HandlerDeleteChirpsByID)
	mux.HandleFunc("POST /api/refresh", apiConfig.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiConfig.HandlerRevoke)
	mux.HandleFunc("POST /admin/reset", apiConfig.HandlerReset)
	mux.HandleFunc("GET /admin/metrics", apiConfig.HandlerGetMetrics)
}
