package routes

import (
	"context"
	"net/http"

	"react-app-golang/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func racesRoutes(e *gin.Engine, client *mongo.Client) {
	rg := e.Group("/races", authMiddleware())

	rg.GET("/", func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		coll, err := client.Database("percent-back-app").Collection("races").Find(context.TODO(), bson.M{"userId": userId})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer coll.Close(context.TODO())

		races := make([]models.Race, 0)
		for coll.Next(context.TODO()) {
			var race models.Race
			if err := coll.Decode(&race); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			races = append(races, race)
		}
		c.JSON(200, races)
	})

	rg.POST("/create", func(c *gin.Context) {
		var race models.CreateRace

		err := c.Bind(&race)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get userId from context, set by authMiddleware
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}
		race.UserID = userId.(primitive.ObjectID) // Assign the authenticated user's ID

		coll := client.Database("percent-back-app").Collection("races")

		doc := models.CreateRace{RaceName: race.RaceName, RaceDate: race.RaceDate, RaceDistance: race.RaceDistance, PercentBack: race.PercentBack, UserID: race.UserID}
		result, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	})

	rg.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid race ID"})
			return
		}

		coll := client.Database("percent-back-app").Collection("races")

		// Verify that the race belongs to the authenticated user
		var race models.Race
		err = coll.FindOne(context.TODO(), bson.M{"_id": objectId, "userId": userId}).Decode(&race)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Race not found or not authorized to delete"})
			return
		}

		result, err := coll.DeleteOne(context.TODO(), bson.M{"_id": objectId})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	})
}

