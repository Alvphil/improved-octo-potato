package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/Alvphil/improved-octo-potato.git/internal/database"
	"github.com/go-chi/chi"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	chirps := []database.Chirp{}

	s := r.URL.Query().Get("author_id")
	// s is a string that contains the value of the author_id query parameter
	// if it exists, or an empty string if it doesn't

	if s != "" {
		s, err := strconv.Atoi(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author ID")
			return
		}
		for _, dbChirp := range dbChirps {
			if dbChirp.Author_id == s {
				chirps = append(chirps, database.Chirp{
					Author_id: dbChirp.Author_id,
					Body:      dbChirp.Body,
					ID:        dbChirp.ID,
				})
			}
		}
	} else {
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, database.Chirp{
				Author_id: dbChirp.Author_id,
				Body:      dbChirp.Body,
				ID:        dbChirp.ID,
			})
		}
	}

	sortedParam := r.URL.Query().Get("sort")
	if sortedParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, database.Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}
