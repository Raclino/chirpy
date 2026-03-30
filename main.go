package main

import "net/http"

func main() {
	muxServer := http.NewServeMux()

	server := &http.Server{
		Handler: muxServer,
		Addr:    ":8080",
	}

	server.ListenAndServe()

}
