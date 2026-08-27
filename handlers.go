package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
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
		respondWithError(w, http.StatusInternalServerError, "Error returning game", err)
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
		Film1ID:   int32(AD.Results[0].KnownFor[0].ID),
		Film2ID:   int32(AD.Results[0].KnownFor[1].ID),
		Film3ID:   int32(AD.Results[0].KnownFor[2].ID),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Game creation failed", err)
		return
	}

	// Create wait group and error chan
	var wait sync.WaitGroup
	errChan := make(chan error, len(AD.Results[0].KnownFor))

	// Loop through films and enter them in database
	for _, film := range AD.Results[0].KnownFor {
		wait.Add(1)

		go func(film KnownFilms) {

			defer wait.Done()

			cfg.enterFilm(film.ID, errChan)

		}(film)
	}

	// Close chan and wait group
	wait.Wait()
	close(errChan)

	// Loop through chan for errors
	for err := range errChan {
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating film in database", err)
			return
		}
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

func (cfg *apiConfig) handlerGameState(w http.ResponseWriter, r *http.Request) {
	type gameState struct {
		Date     string                       `json:"date"`
		Actor    string                       `json:"actor"`
		PlayerID string                       `json:"playerid"`
		Guesses  []database.FetchGuessListRow `json:"guesses"`
		Strikes  int                          `json:"strikes"`
		Posters  []string                     `json:"posters"`
	}

	// Get todays date
	today := time.Now().Format("2006-01-02")

	// Parse player ID
	playerIDString := r.Header.Get("X-Player-ID")
	playerID, err := uuid.Parse(playerIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing UUID", err)
		return
	}
	// Return Game
	game, err := cfg.database.ReturnGame(r.Context(), today)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error returning game", err)
		return
	}
	// Calculate strikes
	strikes, err := cfg.database.StrikeCount(r.Context(), database.StrikeCountParams{
		Date:     today,
		PlayerID: playerID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error calculating strikes", err)
		return
	}
	// Retrieve poster paths from ID
	filmIDList := []int{int(game.Film1ID), int(game.Film2ID), int(game.Film3ID)}
	posterPaths := []string{}

	for _, filmID := range filmIDList {
		path, err := cfg.database.FetchPosters(r.Context(), int32(filmID))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching poster paths", err)
			return
		}
		posterPaths = append(posterPaths, path)
	}

	// Fetch list of guesses
	guesses, err := cfg.database.FetchGuessList(r.Context(), database.FetchGuessListParams{
		Date:     today,
		PlayerID: playerID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, gameState{
				Date:     game.Date,
				Actor:    game.ActorName,
				PlayerID: playerIDString,
				Guesses:  guesses,
				Strikes:  int(strikes),
				Posters:  posterPaths,
			})
			fmt.Println(game.Date, game.ActorName, playerIDString, guesses, strikes)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error fetching guess list", err)
			return
		}
	}

	fmt.Println(game.Date, game.ActorName, playerIDString, guesses, strikes)
	respondWithJSON(w, http.StatusOK, gameState{
		Date:     game.Date,
		Actor:    game.ActorName,
		PlayerID: playerIDString,
		Guesses:  guesses,
		Strikes:  int(strikes),
		Posters:  posterPaths,
	})
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
		respondWithError(w, http.StatusInternalServerError, "Error returning game", err)
		return
	}

	// Retrieve and parse player ID
	playerIDString := r.Header.Get("X-Player-ID")
	playerID, err := uuid.Parse(playerIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing UUID", err)
		return
	}

	// Check if guess matches films, create guess in database
	if strings.EqualFold(params.Guess, game.Film1) {
		payload, err := cfg.guessResponse(params.GameDate, 1, playerID, params.Guess, true)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating guess", err)
			return
		}
		respondWithJSON(w, http.StatusOK, payload)
		return
	}

	if strings.EqualFold(params.Guess, game.Film2) {
		payload, err := cfg.guessResponse(params.GameDate, 2, playerID, params.Guess, true)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating guess", err)
			return
		}
		respondWithJSON(w, http.StatusOK, payload)
		return
	}

	if strings.EqualFold(params.Guess, game.Film3) {
		payload, err := cfg.guessResponse(params.GameDate, 3, playerID, params.Guess, true)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating guess", err)
			return
		}
		respondWithJSON(w, http.StatusOK, payload)
		return
	}

	// Create guess in database for incorrect guess
	payload, err := cfg.guessResponse(params.GameDate, 0, playerID, params.Guess, false)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating guess", err)
		return
	}

	respondWithJSON(w, http.StatusOK, payload)
}
