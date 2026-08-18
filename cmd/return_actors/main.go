package main

import (
	"context"
	"database/sql"

	// "encoding/json"
	"fmt"
	"log"

	// "net/http"
	"os"
	// "strings"
	// "time"

	"github.com/joho/godotenv"

	"GitHub.com/dylandeyotte/movie_guesser/internal/database"

	_ "github.com/lib/pq"
)

// type ActorData struct {
// 	Page    int `json:"page"`
// 	Results []struct {
// 		Adult              bool    `json:"adult"`
// 		Gender             int     `json:"gender"`
// 		ID                 int     `json:"id"`
// 		KnownForDepartment string  `json:"known_for_department"`
// 		Name               string  `json:"name"`
// 		OriginalName       string  `json:"original_name"`
// 		Popularity         float64 `json:"popularity"`
// 		ProfilePath        string  `json:"profile_path"`
// 		KnownFor           []struct {
// 			Adult            bool    `json:"adult"`
// 			BackdropPath     string  `json:"backdrop_path"`
// 			ID               int     `json:"id"`
// 			Title            string  `json:"title"`
// 			OriginalTitle    string  `json:"original_title"`
// 			Overview         string  `json:"overview"`
// 			PosterPath       string  `json:"poster_path"`
// 			MediaType        string  `json:"media_type"`
// 			OriginalLanguage string  `json:"original_language"`
// 			GenreIds         []int   `json:"genre_ids"`
// 			Popularity       float64 `json:"popularity"`
// 			ReleaseDate      string  `json:"release_date"`
// 			Softcore         bool    `json:"softcore"`
// 			Video            bool    `json:"video"`
// 			VoteAverage      float64 `json:"vote_average"`
// 			VoteCount        int     `json:"vote_count"`
// 		} `json:"known_for"`
// 	} `json:"results"`
// 	TotalPages   int `json:"total_pages"`
// 	TotalResults int `json:"total_results"`
// }

func main() {
	// Load data from ENV
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("No DB_URL set")
	}
	// tmdbToken := os.Getenv("TMDB_TOKEN")
	// if tmdbToken == "" {
	// 	log.Fatal("No TMDB token set")
	// }
	// bearerToken := fmt.Sprintf("Bearer %v", tmdbToken)

	// Create db handle
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	actors, err := dbQueries.ReturnActors(context.Background())
	if err != nil {
		fmt.Printf("error returning actors: %v\n", err)
		return
	}

	for _, actor := range actors {
		// client := &http.Client{}

		// splitName := strings.Split(actor.Name, " ")
		// newName := strings.Join(splitName, "+")

		// url := fmt.Sprintf("https://api.themoviedb.org/3/search/person?query=%v", newName)

		// req, err := http.NewRequest("GET", url, nil)
		// if err != nil {
		// 	fmt.Printf("request creation error: %v\n", err)
		// 	return
		// }

		// req.Header.Set("Authorization", bearerToken)

		// resp, err := client.Do(req)
		// if err != nil {
		// 	fmt.Printf("request error: %v\n", err)
		// 	return
		// }

		// defer resp.Body.Close()

		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// bodyString := string(body)

		// fmt.Println(bodyString)
		// var AD ActorData

		// decoder := json.NewDecoder(resp.Body)

		// if err := decoder.Decode(&AD); err != nil {
		// 	fmt.Printf("decoding error: %v\n", err)
		// }
		// fmt.Println(AD)
		// time.Sleep(2 * time.Second)
		// if AD.Results[0].KnownFor[0].MediaType == "tv" || AD.Results[0].KnownFor[1].MediaType == "tv" || AD.Results[0].KnownFor[2].MediaType == "tv" {
		// 	fmt.Printf("%v: TV actor\n", actor.Name)
		// } else {
		// 	fmt.Println(actor.Name)
		// }
		fmt.Println(actor.Name)
	}

	fmt.Println(len(actors))
}
