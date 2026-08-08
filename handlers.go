package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func handlerActorFetch(w http.ResponseWriter, r *http.Request) {

	client := &http.Client{}

	// MUST FETCH ACTOR FROM SOMEWHERE
	actorName := "Robert+Pattinson"

	url := fmt.Sprintf("https://api.themoviedb.org/3/search/person?query=%v&limit=1", actorName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Request creation failed", err)
	}

	godotenv.Load()

	tmdbToken := os.Getenv("TMDB_TOKEN")
	if tmdbToken == "" {
		log.Fatal("No TMDB token set")
	}

	bearerToken := fmt.Sprintf("Bearer %v", tmdbToken)

	req.Header.Set("Authorization", bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Request failed", err)
	}

	defer resp.Body.Close()

	var AD ActorData

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&AD); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
	}

	filmList := []string{}

	for _, film := range AD.Results[0].KnownFor {
		filmList = append(filmList, film.Title)
	}

	type Payload struct {
		Actor string   `json:"actor"`
		Films []string `json:"films"`
	}

	payload := Payload{
		Actor: actorName,
		Films: filmList,
	}

	respondWithJSON(w, http.StatusOK, payload)
}
