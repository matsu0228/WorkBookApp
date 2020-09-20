package main

import (
	"WorkBookApp/internal/api"
	"WorkBookApp/internal/database"
	"context"
	"log"
	"net/http"
)

const PORT = ":8080"

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("api起動完了: http:localhost%s", PORT)
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

	return http.ListenAndServe(PORT, router)
}
