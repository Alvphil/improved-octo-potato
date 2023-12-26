package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerMetricsInc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, cfg.numberServerHits())
}

func (cfg *apiConfig) adminMetricsInc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<html>\n\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n\n	</html>", cfg.fileserverHits)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++ // Increment the counter here
		log.Printf("number of hits: %v", cfg.fileserverHits)
		next.ServeHTTP(w, r) // Call the next handler
	})
}
