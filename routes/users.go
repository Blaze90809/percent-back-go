package routes

import (
	"context"
	"os"
	"react-app-golang/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func usersRoutes(e *gin.Engine) {
	e.POST("/register", func(c *gin.Context) {
		var user models.RegisterUser

		err := c.Bind(&user)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		if user.Username == "" || user.Password == "" {
			c.JSON(401, gin.H{"error": "User needs to enter both a username and a password"})
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
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		coll := client.Database("percent-back-app").Collection("users")

		doc := models.RegisterUser{Username: user.Username, Password: user.Password}
		result, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	})
}
