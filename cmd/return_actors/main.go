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

	actors, err := dbQueries.ReturnActors(context.Background())
	if err != nil {
		fmt.Printf("error returning actors: %v\n", err)
		return
	}

	for _, actor := range actors {
		fmt.Println(actor.Name)
	}
}
