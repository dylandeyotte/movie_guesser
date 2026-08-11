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
		"Yvonne Yung Hung",
		"Yoon Yool",
		"Anastasia Shestakova",
		"Mayuko Sasaki",
		"Ui Mihara",
		"Shoko Fujimura",
		"Julia Doyle",
		"Min Do-yoon",
		"Anaïs Bret",
		"Inde Navarrette",
		"Charlie Cho Cha-Lee",
		"Daniella Wang",
		"Christopher Nolan",
		"Diana Pang Dan",
		"Milly Alcock",
		"Amélie Hoeferle",
		"Sophie Ngan Chin-Man",
		"Ahn Ha-young",
		"Kōichi Imaizumi",
		"Yui Kitamura",
		"Yuka Tano",
		"Mireille Enos",
		"Bridget Regan",
		"Truman Hanks",
		"Jung Hae-in",
		"Riley Chung",
		"Angelica Hart",
		"Simon Baker",
		"Emma Ho",
		"Alan Ritchson",
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
		fmt.Printf("removed %v from database\n", err)
	}
}
