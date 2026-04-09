package main

import (
	"net/http"

	"github.com/Raclino/chirpy/internal/handlers"
)

// Constructor useful for future middleware, see : https://grafana.com/blog/how-i-write-http-services-in-go-after-13-years/#the-new-server-constructor
func NewServer(
	apiConfig *handlers.ApiConfig,
) http.Handler {
	muxServer := http.NewServeMux()
	addRoutes(muxServer, apiConfig)

	var handler http.Handler = muxServer
	return handler
}
