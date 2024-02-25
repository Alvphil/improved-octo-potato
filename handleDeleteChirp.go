package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirpID")
		return
	}
	chirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No valid chirp")
		return
	}

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
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid userID")
		return
	}

	if chirp.Author_id == userID {
		cfg.DB.DeleteChirp(chirpID)
	} else if chirp.Author_id != userID {
		respondWithError(w, http.StatusForbidden, "Invalid token")
		return
	}

}
