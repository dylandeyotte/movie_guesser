package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"GitHub.com/dylandeyotte/movie_guesser/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerActorFetch(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Actor    string `json:"actor"`
		GameDate string `json:"gamedate"`
	}
	// Get todays date
	today := time.Now().Format("2006-01-02")

	// Check if game exists
	game, err := cfg.database.ReturnGame(r.Context(), today)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error returning game, first", err)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// If game exists, return it
		payload := Payload{
			Actor:    game.ActorName,
			GameDate: game.Date,
		}
		respondWithJSON(w, http.StatusOK, payload)
		return
	}

	// Create client
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
	game, err = cfg.database.CreateGame(r.Context(), database.CreateGameParams{
		Date:      today,
		ActorID:   actor.ID,
		ActorName: actor.Name,
		Film1:     AD.Results[0].KnownFor[0].Title,
		Film2:     AD.Results[0].KnownFor[1].Title,
		Film3:     AD.Results[0].KnownFor[2].Title,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Game creation failed", err)
		return
	}

	// if err := cfg.database.MarkActor(r.Context(), actor.Name); err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Error marking actor", err)
	// }

	// Assemble payload
	payload := Payload{
		Actor:    actor.Name,
		GameDate: game.Date,
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerVerifyGuess(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Guess    string `json:"guess"`
		GameDate string `json:"gamedate"`
	}

	params := parameters{}

	// Decode response
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
		return
	}

	// Fetch game from date
	game, err := cfg.database.ReturnGame(r.Context(), params.GameDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error returning game second", err)
		return
	}

	// Retrieve and parse player ID
	id := r.Header.Get("X-Played-ID")
	playerID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing UUID", err)
		return
	}

	type Payload struct {
		Film    string `json:"film"`
		Verdict string `json:"verdict"`
	}

	// Check if guess matches films, create guess in database
	if strings.EqualFold(params.Guess, game.Film1) || strings.EqualFold(params.Guess, game.Film2) || strings.EqualFold(params.Guess, game.Film3) {
		payload, err := cfg.guessResponse(params.GameDate, playerID, params.Guess, true)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating guess", err)
			return
		}
		respondWithJSON(w, http.StatusOK, payload)
		return
	}

	payload, err := cfg.guessResponse(params.GameDate, playerID, params.Guess, false)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating guess", err)
		return
	}

	respondWithJSON(w, http.StatusOK, payload)
}
