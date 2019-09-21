package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/efreitasn/gtodo/cmd/gtodo/internal/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoOptions := options.Client()

	mongoOptions.ApplyURI(os.Getenv("MONGODB_URI"))

	client, err := mongo.Connect(
		ctx,
		mongoOptions,
	)

	if err != nil {
		log.Fatalln(err.Error())
	}

	db := client.Database("gtodo")

	mux := handlers.NewMux(db)
	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	server.ListenAndServeTLS("./.cert/cert.pem", "./.cert/key.pem")
}
