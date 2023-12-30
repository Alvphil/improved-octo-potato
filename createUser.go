package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(params.Email) > 140 {
		w.WriteHeader(400)
		return
	}
	email := params.Email
	passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(params.Password), 0)

	respBody, _ := cfg.DB.CreateUser(email, passwordHashed)

	respondWithJSON(w, http.StatusCreated, respBody)
}
