package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"GitHub.com/dylandeyotte/movie_guesser/internal/database"
)

func (cfg *apiConfig) handlerActorFetch(w http.ResponseWriter, r *http.Request) {

	// Check if game exists

	client := &http.Client{}

	// Select actor from database NEED NEW QUERY NEW ACTOR NOT USED
	actor, err := cfg.database.SelectActor(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching actor", err)
		return
	}
	// Assemble URL
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/person?query=%v&limit=1", actor.Name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Request creation failed", err)
		return
	}
	// Set header token
	req.Header.Set("Authorization", cfg.tmdbToken)

	// HTTP Request
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Request failed", err)
		return
	}
	defer resp.Body.Close()

	var AD ActorData

	// Decode JSON
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&AD); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
		return
	}

	// Create Game
	_, err = cfg.database.CreateGame(r.Context(), database.CreateGameParams{
		Date:    time.Now(),
		ActorID: actor.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Game creation failed", err)
		return
	}

	// Assemble film list
	filmList := []string{}
	for _, film := range AD.Results[0].KnownFor {
		filmList = append(filmList, film.Title)
	}

	// if err := cfg.database.MarkActor(r.Context(), actor.Name); err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Error marking actor", err)
	// }

	// Assemble payload
	type Payload struct {
		Actor string `json:"actor"`
	}

	payload := Payload{
		Actor: actor.Name,
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerVerifyGuess(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Guess string `json:"guess"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
		return
	}

	// Pull answer from somewhere
	correctGuess := "Batman"

	if strings.EqualFold(params.Guess, correctGuess) {
		respondWithJSON(w, http.StatusOK, "Correct")
		return
	}

	// Store guess in database

	respondWithJSON(w, http.StatusOK, "Incorrect")
}
