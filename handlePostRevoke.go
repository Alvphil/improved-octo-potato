package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
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

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}
	if claims.Issuer == "chirpy-refresh" {
		cfg.DB.RevokeRefreshToken(token.Raw)
		w.WriteHeader(200)
		return
	}

}
