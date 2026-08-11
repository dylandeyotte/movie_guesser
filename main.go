package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/joho/godotenv"

	"GitHub.com/dylandeyotte/movie_guesser/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	database  *database.Queries
	tmdbToken string
}

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("No DB_URL set")
	}
	tmdbToken := os.Getenv("TMDB_TOKEN")
	if tmdbToken == "" {
		log.Fatal("No TMDB token set")
	}
	bearerToken := fmt.Sprintf("Bearer %v", tmdbToken)

	// Create db handle
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	apiCfg := apiConfig{
		database:  dbQueries,
		tmdbToken: bearerToken,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/actor", apiCfg.handlerActorFetch)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	defer server.Close()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
