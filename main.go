package main

import (
	"WorkBookApp/internal/api"
	"WorkBookApp/internal/database"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	dclient, sclient, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dclient.Close(); err != nil {
			log.Println("can't close", err)
		}
		if err := sclient.Close(); err != nil {
			log.Println("can't close", err)
		}
	}()

	log.Printf("api起動完了: http:localhost%s", api.PORT)
}

func run() (*datastore.Client, *storage.Client, error) {
	ctx := context.Background()
	DataStoreClient, err := database.DataStoreNewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	StoregaeClient, err := api.CloudStoreNewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	router := api.Route()
	log.Printf("Listening on port %s", port)
	return DataStoreClient, StoregaeClient, http.ListenAndServe(":"+port, router)
}
