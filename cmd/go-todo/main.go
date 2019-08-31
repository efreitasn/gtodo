package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/efreitasn/go-todo/cmd/go-todo/internal/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoOptions := options.Client()

	mongoOptions.ApplyURI(os.Getenv("MONGODB_URL"))
	mongoOptions.SetAuth(options.Credential{
		Username:    os.Getenv("MONGODB_USERNAME"),
		Password:    os.Getenv("MONGODB_PASSWORD"),
		PasswordSet: os.Getenv("MONGODB_PASSWORD") == "",
	})

	client, err := mongo.Connect(
		ctx,
		mongoOptions,
	)

	if err != nil {
		log.Fatalln(err.Error())
	}

	db := client.Database(os.Getenv("MONGODB_DB"))

	mux := handlers.NewMux(db)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServeTLS("./cert/cert.pem", "./cert/key.pem")
}
