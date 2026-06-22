package main

import (
	"log"

	"github.com/omzamirr/internal/server"
	"github.com/omzamirr/internal/store"
)

func main() {
	kvStore := store.New()

	srv := server.New(":6380", kvStore)

	log.Println("Attempting to bind to port :6380...")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Port binding failed: %v\n", err)
	}

	log.Println("Success! Port bound and released cleanly with no errors.")
}
