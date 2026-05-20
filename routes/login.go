package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"react-app-golang/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func loginRoutes(e *gin.RouterGroup, client *mongo.Client) {
	e.POST("/login", RateLimitMiddleware(SensitiveLimiter), func(c *gin.Context) {
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

		var authUser models.User
		err = client.Database("percent-back-app").Collection("users").FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&authUser)
		if err != nil {
			fmt.Println("User lookup failed:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(authUser.Password), []byte(user.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
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
			fmt.Println("Failed to sign JWT token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": signedToken})
	})
}
