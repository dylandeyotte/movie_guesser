package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"GitHub.com/dylandeyotte/movie_guesser/internal/database"
	"github.com/google/uuid"
)

type Payload struct {
	Verdict    bool     `json:"verdict"`
	FilmNumber int      `json:"film_number"`
	Strikes    int      `json:"strikes"`
	Guess      string   `json:"guess"`
	PlayerID   string   `json:"playerid"`
	Repeat     bool     `json:"repeat"`
	PosterPath []string `json:"poster_path"`
	GameOver   bool     `json:"game_over"`
	Victory    bool     `json:"victory"`
}

func (cfg *apiConfig) answerReveal(strikes int, game database.Game) (FilmAnswers, error) {
	if strikes < 3 {
		return FilmAnswers{}, nil
	}
	paths, err := cfg.fetchPosterPaths([]int{int(game.Film1ID), int(game.Film2ID), int(game.Film3ID)})
	if err != nil {
		return FilmAnswers{}, err
	}
	return FilmAnswers{
		game.Film1,
		paths[0],
		game.Film2,
		paths[1],
		game.Film3,
		paths[2],
	}, nil
}

func gameOverCheck(strikes int) bool {
	return strikes >= 3
}

func correctGuessesCheck(guesses []database.Guess) int {
	count := 0
	for _, guess := range guesses {
		if guess.Verdict == true {
			count++
		}
	}
	return count
}

func (cfg *apiConfig) actorFetchHelper() (ActorData, database.Actor, error) {
	for i := range 5 {
		i++
		// Create client
		client := &http.Client{}

		// Select actor from database
		actor, err := cfg.database.SelectActor(context.Background())
		if err != nil {
			return ActorData{}, database.Actor{}, err
		}
		// Assemble URL
		url := fmt.Sprintf("https://api.themoviedb.org/3/search/person?query=%v&limit=1", actor.Name)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ActorData{}, database.Actor{}, err
		}
		// Set header token
		req.Header.Set("Authorization", cfg.tmdbToken)

		// HTTP Request
		resp, err := client.Do(req)
		if err != nil {
			return ActorData{}, database.Actor{}, err
		}
		defer resp.Body.Close()

		var AD ActorData

		// Decode JSON
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&AD); err != nil {
			return ActorData{}, database.Actor{}, err
		}

		if AD.Results[0].KnownFor[0].MediaType == "tv" || AD.Results[0].KnownFor[1].MediaType == "tv" || AD.Results[0].KnownFor[2].MediaType == "tv" {
			continue
		}
		return AD, actor, nil
	}
	return ActorData{}, database.Actor{}, errors.New("Failed to find actor")
}

func (cfg *apiConfig) createUserGameHelper(playerID uuid.UUID, date, actor string, correctGuesses, incorrectGuesses int, victory bool) error {
	_, err := cfg.database.CreateUserGame(context.Background(), database.CreateUserGameParams{
		Date:             date,
		PlayerID:         playerID,
		Actor:            actor,
		CorrectGuesses:   int32(correctGuesses),
		IncorrectGuesses: int32(incorrectGuesses),
		Victory:          victory,
	})
	if err != nil {
		return err
	}
	return nil
}

func (cfg *apiConfig) fetchPosterPaths(filmIDList []int) ([]string, error) {
	posterPaths := []string{}

	for _, filmID := range filmIDList {
		path, err := cfg.database.FetchPosters(context.Background(), int32(filmID))
		if err != nil {
			return []string{}, err
		}
		posterPaths = append(posterPaths, path)
	}
	return posterPaths, nil
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

func (cfg *apiConfig) guessResponse(date string, filmNumber, filmID int, playerID uuid.UUID, guess, actor string, verdict bool) (Payload, error) {
	// If guess is correct, fetch poster path
	posterPath := []string{}
	if verdict == true {
		poster, err := cfg.fetchPosterPaths([]int{filmID})
		posterPath = append(posterPath, poster[0])
		if err != nil {
			fmt.Println(err)
			return Payload{}, err
		}
	}
	// Get list of correct guesses to check for victory
	// This is prior to the guess entering the db in case of dulpicate guess
	guesses, err := cfg.database.FetchGuessList(context.Background(), database.FetchGuessListParams{
		Date:     date,
		PlayerID: playerID,
	})
	if err != nil {
		fmt.Println(err)
		return Payload{}, err
	}
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
			fmt.Println(err)
			return Payload{}, err
		}
		fmt.Printf("%v repeated the guess: %v\n", playerID, guess)
		return Payload{
			Verdict:    fetchedGuess.Verdict,
			FilmNumber: int(fetchedGuess.FilmNumber),
			Strikes:    int(strikes),
			Guess:      fetchedGuess.Guess,
			PlayerID:   uuid.UUID.String(fetchedGuess.PlayerID),
			Repeat:     true,
			PosterPath: posterPath,
			GameOver:   gameOverCheck(int(strikes)),
			Victory:    correctGuessesCheck(guesses) == 3,
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
		fmt.Println(err)
		return Payload{}, err
	}
	// Get list of correct guesses to check for victory
	// This is after the guess is entered into the db, the updated list
	updatedGuesses, err := cfg.database.FetchGuessList(context.Background(), database.FetchGuessListParams{
		Date:     date,
		PlayerID: playerID,
	})
	if err != nil {
		fmt.Println(err)
		return Payload{}, err
	}
	// Calculate strike count
	strikes, err := cfg.database.StrikeCount(context.Background(), database.StrikeCountParams{
		Date:     date,
		PlayerID: playerID,
	})
	if err != nil {
		fmt.Println(err)
		return Payload{}, err
	}
	// If victory, create completed game in db
	if correctGuessesCheck(updatedGuesses) == 3 {
		if err := cfg.createUserGameHelper(playerID, date, actor, 3, int(strikes), true); err != nil {
			fmt.Println(err)
			return Payload{}, err
		}
	}
	// If defeat, create completed game in db
	if gameOverCheck(int(strikes)) {
		if err := cfg.createUserGameHelper(playerID, date, actor, correctGuessesCheck(updatedGuesses), 3, false); err != nil {
			fmt.Println(err)
			return Payload{}, err
		}
	}

	// Return payload
	return Payload{
		Verdict:    verdict,
		FilmNumber: filmNumber,
		Strikes:    int(strikes),
		Guess:      guess,
		PlayerID:   uuid.UUID.String(playerID),
		Repeat:     false,
		PosterPath: posterPath,
		GameOver:   gameOverCheck(int(strikes)),
		Victory:    correctGuessesCheck(updatedGuesses) == 3,
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
