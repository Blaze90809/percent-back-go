package routes

import (
	"context"
	"react-app-golang/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func usersRoutes(e *gin.Engine, client *mongo.Client) {
	e.POST("/register", func(c *gin.Context) {
		var user models.RegisterUser

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		if user.Username == "" || user.Password == "" {
			c.JSON(401, gin.H{"error": "User needs to enter both a username and a password"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		coll := client.Database("percent-back-app").Collection("users")

		doc := models.RegisterUser{Username: user.Username, Password: string(hashedPassword)}
		result, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	})
}
