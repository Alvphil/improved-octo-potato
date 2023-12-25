package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	r := chi.NewRouter()
	apiCfg := apiConfig{0}
	r.Mount("/api/", api(&apiCfg))

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)
	//r.Get("/healthz", handlerReadiness)
	//r.Get("/metrics", apiCfg.handlerMetricsInc)
	//r.HandleFunc("/reset", apiCfg.handlerResetMetrics)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetricsInc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, cfg.numberServerHits())
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}

func (cfg apiConfig) numberServerHits() string {
	return ("Hits: " + fmt.Sprint(cfg.fileserverHits))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++ // Increment the counter here
		log.Printf("number of hits: %v", cfg.fileserverHits)
		next.ServeHTTP(w, r) // Call the next handler
	})

	//cfg.fileserverHits++
	//log.Printf("number of hits: %v", cfg.fileserverHits)
	//return next
}

func api(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", apiCfg.handlerMetricsInc)
	r.HandleFunc("/reset", apiCfg.handlerResetMetrics)
	return r
}
