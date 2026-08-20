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
	Verdict bool `json:"verdict"`
	Strikes int  `json:"strikes"`
}

// ALREADY GUESSED???
func (cfg *apiConfig) guessResponse(date string, playerID uuid.UUID, guess string, verdict bool) (Payload, error) {
	_, err := cfg.database.CreateGuess(context.Background(), database.CreateGuessParams{
		Date:     date,
		PlayerID: playerID,
		Guess:    guess,
		Verdict:  verdict,
	})
	if err != nil {
		return Payload{}, err
	}

	strikes, err := cfg.database.StrikeCount(context.Background(), database.StrikeCountParams{
		Date:     date,
		PlayerID: playerID,
	})
	if err != nil {
		return Payload{}, err
	}
	return Payload{
		Verdict: verdict,
		Strikes: int(strikes),
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
