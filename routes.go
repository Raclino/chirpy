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
	mux.HandleFunc("GET /api/healthz", handlers.HandleGetHealth)
	mux.HandleFunc("POST /api/users", apiConfig.HandleCreateUsers)
	mux.HandleFunc("PUT /api/users", apiConfig.HandleUpdateUsers)
	mux.HandleFunc("POST /api/login", apiConfig.HandleUserLogin)
	mux.HandleFunc("POST /api/chirps", apiConfig.HandleCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiConfig.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.HandleGetChirpsByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.HandleDeleteChirpsByID)
	mux.HandleFunc("POST /api/refresh", apiConfig.HandleRefresh)
	mux.HandleFunc("POST /api/revoke", apiConfig.HandleRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.HandlePolkaWebhooks)
	mux.HandleFunc("POST /admin/reset", apiConfig.HandleReset)
	mux.HandleFunc("GET /admin/metrics", apiConfig.HandleGetMetrics)
}
