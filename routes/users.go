package routes

import (
	"context"
	"fmt"
	"net/http"
	"react-app-golang/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func usersRoutes(e *gin.RouterGroup, client *mongo.Client) {
	e.POST("/register", RateLimitMiddleware(SensitiveLimiter), func(c *gin.Context) {
		var user models.RegisterUser

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if user.Username == "" || user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
			return
		}

		// Enforce password strength minimum length (8 characters)
		if len(user.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
			return
		}

		coll := client.Database("percent-back-app").Collection("users")

		// Programmatically verify username uniqueness before registration
		var existingUser models.User
		err = coll.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Failed to hash password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		doc := models.RegisterUser{Username: user.Username, Password: string(hashedPassword)}
		result, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			fmt.Println("Failed to insert user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, result)
	})
}
