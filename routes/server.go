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

type CreateRace struct {
    RaceName     string `json:"raceName"`
    RaceDate     string `json:"raceDate"`
    RaceDistance float64 `json:"raceDistance"`
    PercentBack  float64 `json:"percentBack"`
    UserID       string `json:"userId"`
}

func NewRouter() {
	e := gin.Default()

	e.GET("/", func(c *gin.Context) {
		c.String(200, "Nordic percent back.")
	})

	e.GET("/races/:objectId", func(c *gin.Context) {
		id := c.Param("objectId")

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

		objectId, err := primitive.ObjectIDFromHex(id)
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

	e.POST("/races/create", func(c *gin.Context) {
		var race CreateRace

		err := c.Bind(&race)
		if err != nil {
			panic(err)
		}

		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		err = godotenv.Load()
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
	
		coll := client.Database("percent-back-app").Collection("races")
		if err != nil{
    		panic(err)
		}
		doc := CreateRace{RaceName: race.RaceName, RaceDate: race.RaceDate, RaceDistance: race.RaceDistance, PercentBack: race.PercentBack, UserID: race.UserID}
		result, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			panic(err)
		}

		c.JSON(200, result)
	})

	err := e.Run()
	if err != nil {
		panic(err)
	}
}