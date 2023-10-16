package main

import (
	"context"
	"fmt"
	"os"
	server "react-app-golang/routes"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	connectionURI := os.Getenv("mongo_uri")

	opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
	  panic(err)
	}

	defer func() {
	  if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	  }
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	  panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	server.NewRouter()
}