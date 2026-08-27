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
		"Dawn Olivieri",
		"Dax Flame",
		"Olivier Richters",
		"Ismail Bashey",
		"Anil Kapoor",
		"Fan Bingbing",
		"Genevieve Buechner",
		"Shu Qi",
		"Christopher Rodriguez Marquette",
		"Bo Dietl",
		"Yuko Tanaka",
		"Kyle Chandler",
		"Emmanuelle Chriqui",
		"Kelly Reilly",
		"Elisabeth Röhm",
		"Peter Stormare",
		"Ali Larter",
		"Famke Janssen",
		"Max Martini",
		"Kevin Corrigan",
		"Neal McDonough",
		"Claire Forlani",
		"Lochlyn Munro",
		"Brian Tee",
		"Marc Blucas",
		"Laura Haddock",
		"Rebecca Romijn",
		"Gabriel Guevara",
		"Lisa Ann Walter",
		"Keiko Kitagawa",
		"Shawn Ashmore",
		"Callum Turner",
		"Dolph Lundgren",
		"Scott Adkins",
		"Frances Fisher",
		"Jean Reno",
		"Colin Hanks",
		"Olga Kurylenko",
		"Michele Morrone",
		"Sam Claflin",
		"Gong Li",
		"Dermot Mulroney",
		"Bruce McGill",
		"Maria Bello",
		"Sam Elliott",
		"Hope Davis",
		"Ana Ularu",
		"Jonathan Rhys Meyers",
		"Tracey Ullman",
		"Yu Nan",
		"Mireille Enos",
		"Kaya Scodelario",
		"Malcolm McDowell",
		"Common",
		"Cho Yeo-jeong",
		"Stephen Graham",
		"Rhona Mitra",
		"Macaulay Culkin",
		"Don Lee",
		"Laura Linney",
		"Cole Hauser",
		"Alanna Ubach",
		"Élodie Yung",
		"Sophie Marceau",
		"Zhu Zhu",
		"Daniel Wu",
		"Joely Richardson",
		"Liu Yifei",
		"Ingrid Bolsø Berdal",
		"Téa Leoni",
		"Marlon Wayans",
		"Lynn Collins",
		"Hiroyuki Sanada",
		"Giovanni Ribisi",
		"Hana Takeda",
		"Luis Guzmán",
		"Carrie Coon",
		"Cameron Monaghan",
		"Choi Woo-shik",
		"Tom Felton",
		"Krista Allen",
		"Minka Kelly",
		"Richard Armitage",
		"Richard T. Jones",
		"Ian McShane",
		"David Morse",
		"Topher Grace",
		"Elizabeth McGovern",
		"Ed Skrein",
		"Jing Tian",
		"Scott Eastwood",
		"Josh Stewart",
		"Christopher McDonald",
		"Daniel Dae Kim",
		"Leslie Nielsen",
		"Ruby Rose",
		"Samantha Mathis",
		"Gary Cole",
		"Sasha Calle",
		"Amanda Peet",
		"Bob Hoskins",
		"Ming-Na Wen",
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
