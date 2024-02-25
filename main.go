package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Alvphil/improved-octo-potato.git/internal/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	godotenv.Load()
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}
	r := chi.NewRouter()

	r.Mount("/api/", api(&apiCfg))
	r.Mount("/admin/", admin(&apiCfg))

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

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

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}

func (cfg apiConfig) numberServerHits() string {
	return ("Hits: " + fmt.Sprint(cfg.fileserverHits))
}

func (cfg apiConfig) adminNumberServerHits() string {
	return (fmt.Sprint(cfg.fileserverHits))
}

func api(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", apiCfg.handlerMetricsInc)
	r.Post("/chirps", apiCfg.handlerPostChirp)
	r.Get("/chirps", apiCfg.handlerGetChirps)
	r.Get("/chirps/{chirpID}", apiCfg.handlerGetChirp)
	r.Post("/users", apiCfg.handlerPostUser)
	r.Put("/users", apiCfg.HandlerPutUsers)
	r.Get("/users/{userID}", apiCfg.handlerGetUser)
	r.Post("/login", apiCfg.HandlerLoginUser)
	r.Post("/refresh", apiCfg.HandlerRefreshToken)
	r.Post("/revoke", apiCfg.HandlerRevokeToken)
	r.HandleFunc("/reset", apiCfg.handlerResetMetrics)
	return r
}

func admin(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()
	r.Get("/metrics", apiCfg.adminMetricsInc)
	return r
}
