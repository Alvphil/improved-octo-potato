package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
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
	if claims.Issuer == "chirpy-refresh" {
		cfg.DB.RevokeRefreshToken(token.Raw)
		w.WriteHeader(200)
		return
	}

}
