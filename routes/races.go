package routes

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"react-app-golang/models"
)

func racesRoutes(e *gin.Engine) {
	e.GET("/races/:objectId", func(c *gin.Context) {
		id := c.Param("objectId")
	
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		err := godotenv.Load()
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		connectionURI := os.Getenv("mongo_uri")
	
		opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)
		// Create a new client and connect to the server
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		coll, err := client.Database("percent-back-app").Collection("races").Find(context.TODO(), bson.M{"userId": objectId})
		if err != nil {
			panic(err)
		}
		defer coll.Close(context.TODO())
	
		races := make([]models.Race, 0)
		for coll.Next(context.TODO()) {
			var race models.Race
			if err := coll.Decode(&race); err != nil {
				c.JSON(401, gin.H{"error": err.Error()})
				return
			}
			races = append(races, race)
		}
		c.JSON(200, races)
	})
	
	e.POST("/races/create", func(c *gin.Context) {
		var race models.CreateRace
	
		err := c.Bind(&race)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		err = godotenv.Load()
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		connectionURI := os.Getenv("mongo_uri")
	
		opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)
		// Create a new client and connect to the server
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		coll := client.Database("percent-back-app").Collection("races")
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		doc := models.CreateRace{RaceName: race.RaceName, RaceDate: race.RaceDate, RaceDistance: race.RaceDistance, PercentBack: race.PercentBack, UserID: race.UserID}
		result, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(200, result)
	})
	
	e.DELETE("/races/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
	
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		err := godotenv.Load()
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		connectionURI := os.Getenv("mongo_uri")
	
		opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)
		// Create a new client and connect to the server
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		coll := client.Database("percent-back-app").Collection("races")
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		result, err := coll.DeleteOne(context.TODO(), bson.M{"_id": objectId})
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(200, result)
	})
}
