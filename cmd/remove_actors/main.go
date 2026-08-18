package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"GitHub.com/dylandeyotte/movie_guesser/internal/database"

	_ "github.com/lib/pq"
)

func main() {

	removeList := []string{
		"Minka Kelly",
	}
	// Load data from ENV
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("No DB_URL set")
	}
	// Create db handle
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	for _, actor := range removeList {
		if err := dbQueries.RemoveActor(context.Background(), actor); err != nil {
			fmt.Printf("Failed to remove actor %v from database: %v\n", actor, err)
		}
		fmt.Printf("removed %v from database\n", actor)
	}
}
