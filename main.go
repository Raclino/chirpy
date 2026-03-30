package main

import "net/http"

func main() {
	const filePathRoot = "."
	const port = "8080"
	muxServer := http.NewServeMux()
	muxServer.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot))))
	muxServer.Handle("/assets/logo.png", http.StripPrefix("/assets/logo.png", http.FileServer(http.Dir(filePathRoot))))
	muxServer.HandleFunc("/healthz", handlerHealth)

	server := &http.Server{
		Handler: muxServer,
		Addr:    ":" + port,
	}
	server.ListenAndServe()

}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}
