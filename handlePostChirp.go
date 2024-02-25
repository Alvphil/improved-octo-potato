package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	token, err := cfg.extractValidateToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(params.Body) > 140 {
		w.WriteHeader(400)
		return
	}

	//Here is what is going to be returned from a post request:
	cleanedString := cleanBody(params.Body)

	respBody, _ := cfg.DB.CreateChirp(cleanedString, claims.Subject)

	respondWithJSON(w, http.StatusCreated, respBody)

}

func cleanBody(body string) string {
	splitBody := strings.Split(body, " ")
	dissallowedWords := []string{"kerfuffle", "sharbert", "fornax"}
	var cleanB []string
	for _, word := range splitBody {
		if contains(dissallowedWords, strings.ToLower(word)) {
			word = "****"
		}
		cleanB = append(cleanB, word)
	}
	cleanString := strings.Join(cleanB, " ")
	return cleanString
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
