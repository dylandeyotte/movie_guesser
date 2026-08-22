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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Player-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
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
	mux.HandleFunc("POST /api/guess", apiCfg.handlerVerifyGuess)
	mux.HandleFunc("GET /api/gamestate", apiCfg.handlerGameState)

	server := http.Server{
		Handler: corsMiddleware(mux),
		Addr:    ":8080",
	}

	defer server.Close()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
