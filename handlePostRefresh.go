package main

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshedToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
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
	revoked, err := cfg.DB.GetRevokedToken(token.Raw)
	if err == nil {
		if claims.Issuer != "chirpy-refresh" || revoked == token.Raw || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error when finding user ID")
		return
	}

	new_token, err := cfg.CreateJWTToken(3600, id, "chirpy-access")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error when finding user ID")
		return
	}

	RefreshedToken := RefreshedToken{
		Token: new_token,
	}
	respondWithJSON(w, http.StatusOK, RefreshedToken)
}
