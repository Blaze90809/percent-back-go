package server

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Race struct {
    ID          primitive.ObjectID `bson:"_id"`
    RaceName    string             `bson:"raceName"`
    RaceDate    string          `bson:"raceDate"`
    RaceDistance float64                `bson:"raceDistance"`
    PercentBack float64            `bson:"percentBack"`
}

func NewRouter() {
	e := gin.Default()

	e.GET("/", func(c *gin.Context) {
		c.String(200, "Nordic percent back.")
	})

	e.GET("/races", func(c *gin.Context) {
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

		objectId, err := primitive.ObjectIDFromHex("5b4eca874d01b900243f0da7")
		if err != nil{
    		panic(err)
		}

		coll, err := client.Database("percent-back-app").Collection("races").Find(context.TODO(), bson.M{"userId": objectId})
		if err != nil{
    		panic(err)
		}
		defer coll.Close(context.TODO())

		races := make([]Race, 0)
		for coll.Next(context.TODO()) {
			var race Race
			if err := coll.Decode(&race); err != nil {
				panic(err)
			}
			races = append(races, race)
		}
		c.JSON(200, races)
	})

	err := e.Run()
	if err != nil {
		panic(err)
	}
}