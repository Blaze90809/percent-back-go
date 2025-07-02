package routes

import (
	"context"
	"os"
	"react-app-golang/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func loginRoutes(e *gin.Engine, client *mongo.Client) {
	e.POST("/login", func(c *gin.Context) {
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

		var authUser models.User
		err = client.Database("percent-back-app").Collection("users").FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&authUser)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(authUser.Password), []byte(user.Password))
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		claims := models.Claims{
			UserID: authUser.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "percent-back-app",
				Subject:   authUser.Username,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		secretKey := os.Getenv("JWT_SECRET_KEY")
		signedToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to sign token"})
			return
		}

		c.JSON(200, gin.H{"token": signedToken})
	})
}
