package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func (cfg *apiConfig) HandlerChirpyRed(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	words := strings.Split(authHeader, " ")
	if len(words) != 2 {
		respondWithError(w, http.StatusUnauthorized, "User not found")
		return
	}
	PolkaApiKey := os.Getenv("POLKA_API_KEY")

	apiKey := words[1]
	if apiKey != PolkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "User not found")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if params.Event == "user.upgraded" {
		err := cfg.DB.ApplyChirpRed(params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "User not found")
		}
		respondWithJSON(w, http.StatusOK, nil)
	}

}
