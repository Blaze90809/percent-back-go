package routes

import (
	"context"
	"os"
	"react-app-golang/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"time"
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

		connectionURI := os.Getenv("mongo_uri")
		serverApi := options.ServerAPI(options.ServerAPIVersion1)
		err = godotenv.Load()
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverApi)

		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		var authUser models.User
		err = client.Database("percent-back-app").Collection("users").FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&authUser)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		if user.Username != authUser.Username {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(authUser.Password), []byte(user.Password))
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		claims := models.Claims{
			UserID: authUser.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: "percent-back-app",
				Subject: authUser.Username,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
			},
		}
	
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
		// secretKey := []byte(os.Getenv("JWT_SECRET_KEY")) Todo: Implement this.
		secretKey := "test123" // Replace with your actual secret key
		signedToken, err := token.SignedString([]byte(secretKey))
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, signedToken)
	})

}