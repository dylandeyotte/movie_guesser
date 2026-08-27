package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"GitHub.com/dylandeyotte/movie_guesser/internal/database"
	"github.com/google/uuid"
)

type Payload struct {
	Verdict    bool   `json:"verdict"`
	FilmNumber int    `json:"film_number"`
	Strikes    int    `json:"strikes"`
	Guess      string `json:"guess"`
	PlayerID   string `json:"playerid"`
	Repeat     bool   `json:"repeat"`
}

func (cfg *apiConfig) enterFilm(movieID int, errCh chan error) {
	// Create client
	client := &http.Client{}

	// Assemble URL
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%v", movieID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errCh <- err
		return
	}
	// Set header token
	req.Header.Set("Authorization", cfg.tmdbToken)

	// HTTP Request
	resp, err := client.Do(req)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	var FD FilmData

	// Decode JSON
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&FD); err != nil {
		errCh <- err
		return
	}
	// Enter film in database
	if err := cfg.database.CreateFilm(context.Background(), database.CreateFilmParams{
		ID:         int32(FD.ID),
		Title:      FD.Title,
		PosterPath: FD.PosterPath,
	}); err != nil {
		errCh <- err
		return
	}

}

func (cfg *apiConfig) guessResponse(date string, filmNumber int, playerID uuid.UUID, guess string, verdict bool) (Payload, error) {

	// Check if guess has been guessed and return if so
	fetchedGuess, err := cfg.database.FetchGuess(context.Background(), database.FetchGuessParams{
		Date:     date,
		PlayerID: playerID,
		Guess:    guess,
	})
	// If film has been guessed, get strikes and return
	if err == nil {
		// Calculate strike count
		strikes, err := cfg.database.StrikeCount(context.Background(), database.StrikeCountParams{
			Date:     date,
			PlayerID: playerID,
		})
		if err != nil {
			return Payload{}, err
		}
		return Payload{
			Verdict:    fetchedGuess.Verdict,
			FilmNumber: int(fetchedGuess.FilmNumber),
			Strikes:    int(strikes),
			Guess:      fetchedGuess.Guess,
			PlayerID:   uuid.UUID.String(fetchedGuess.PlayerID),
			Repeat:     true,
		}, nil
	}
	// Create guess in database if new
	_, err = cfg.database.CreateGuess(context.Background(), database.CreateGuessParams{
		Date:       date,
		FilmNumber: int32(filmNumber),
		PlayerID:   playerID,
		Guess:      guess,
		Verdict:    verdict,
	})
	if err != nil {
		return Payload{}, err
	}
	// Calculate strike count
	strikes, err := cfg.database.StrikeCount(context.Background(), database.StrikeCountParams{
		Date:     date,
		PlayerID: playerID,
	})
	if err != nil {
		return Payload{}, err
	}
	// Return payload
	return Payload{
		Verdict:    verdict,
		FilmNumber: filmNumber,
		Strikes:    int(strikes),
		Guess:      guess,
		PlayerID:   uuid.UUID.String(playerID),
		Repeat:     false,
	}, nil
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
		fmt.Println(msg)
	}
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Marshal data to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Marshalling error", err)
		return
	}
	// Write response
	w.WriteHeader(code)
	w.Write(data)
}
