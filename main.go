package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/actor", handlerActorFetch)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	defer server.Close()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
