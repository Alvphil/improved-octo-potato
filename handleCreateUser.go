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
	email := params.Email
	passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(params.Password), 0)

	respBody, err := cfg.DB.CreateUser(email, passwordHashed)
	if err != nil {
		respondWithError(w, http.StatusNotAcceptable, "User already exists")
	} else {
		respondWithJSON(w, http.StatusCreated, respBody)
	}

}
