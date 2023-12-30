package main

import (
	"net/http"
	"strconv"

	"github.com/Alvphil/improved-octo-potato.git/internal/database"
	"github.com/go-chi/chi"
)

// func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
// 	dbChirps, err := cfg.DB.GetChirps()
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user")
// 		return
// 	}
// 	chirps := []database.Chirp{}
// 	for _, dbChirp := range dbChirps {
// 		chirps = append(chirps, database.Chirp{
// 			ID:   dbChirp.ID,
// 			Body: dbChirp.Body,
// 		})
// 	}

// 	sort.Slice(chirps, func(i, j int) bool {
// 		return chirps[i].ID < chirps[j].ID
// 	})

// 	respondWithJSON(w, http.StatusOK, chirps)
// }

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	dbUser, err := cfg.DB.GetUser(userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user")
		return
	}

	respondWithJSON(w, http.StatusOK, database.User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
	})
}
