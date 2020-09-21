package main

import (
	"WorkBookApp/internal/api"
	"WorkBookApp/internal/database"
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	//
	d, err := database.NewClient(ctx)
	if err != nil {
		return err
	}

	//
	s, err := api.NewClient(ctx)
	if err != nil {
		return err
	}

	//
	app := api.NewApp(d, s)
	router := api.Route(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	return http.ListenAndServe(":"+port, router)
}
