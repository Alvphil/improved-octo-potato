package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		//ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
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
	password := params.Password

	user, err := cfg.DB.CheckUserExists(email)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHashed, []byte(password)); err != nil {
		// If password does not match
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := cfg.CreateJWTToken(3600, user.ID, "chirpy-access") // 1 hour in seconds = 3600
	if err != nil {
		// If there is an error fetching the user
		respondWithError(w, http.StatusInternalServerError, "Server could not create JWS token for the user")
		return
	}

	refresh, err := cfg.CreateJWTToken(5184000, user.ID, "chirpy-refresh") // 60 days in seconds = 5184000
	if err != nil {
		// If there is an error fetching the user
		respondWithError(w, http.StatusInternalServerError, "Server could not create refresh token for the user")
		return
	}

	resp, err := cfg.DB.GetLoggedInUser(user.ID, token, refresh)
	if err != nil {
		// If there is an error fetching the user
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) CreateJWTToken(expiresInSeconds int, userID int, issuer string) (signedtoken string, err error) {
	now := time.Now().UTC()
	expirationTime := now.Add(time.Duration(expiresInSeconds) * time.Second)

	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   strconv.Itoa(userID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(cfg.jwtSecret)
	signedtoken, err = token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		return "", err
	}
	return signedtoken, nil
}
