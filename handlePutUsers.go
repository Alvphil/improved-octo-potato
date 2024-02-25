package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Alvphil/improved-octo-potato.git/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type UserUpdate struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (cfg *apiConfig) HandlerPutUsers(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	words := strings.Split(authHeader, " ")

	if len(words) != 2 {
		respondWithError(w, http.StatusBadRequest, "Missing auth parameter")
		return
	}

	tokenString := words[1]

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid issuer")
		return
	}
	if issuer == "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var update UserUpdate
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&update)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid userID")
		return
	}

	dbPublicUser, err := cfg.DB.UpdateUser(userID, update.Email, update.Password)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, database.PublicUser{
		ID:    dbPublicUser.ID,
		Email: dbPublicUser.Email,
	})
}
