package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/efreitasn/go-todo/cmd/go-todo/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("MONGODB_URL")),
	)

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")

		os.Exit(1)
	}

	db := client.Database(os.Getenv("MONGODB_DB"))

	mux := handlers.Setup(db)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
