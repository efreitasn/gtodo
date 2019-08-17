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

	client, _ := mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("MONGODB_URL")),
	)
	db := client.Database("go-todo")

	mux := handlers.Setup(db)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
