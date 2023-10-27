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
)

func usersRoutes(e *gin.Engine) {
	e.POST("/register", func(c *gin.Context) {
		var user models.RegisterUser

		err := c.Bind(&user)
		if err != nil {
			panic(err)
		}

		if user.Username == "" || user.Password == "" {
			panic("User needs to enter both a username and a password")
		}

		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		err = godotenv.Load()
		if err != nil {
			panic(err)
		}
		connectionURI := os.Getenv("mongo_uri")

		opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			panic(err)
		}

		coll := client.Database("percent-back-app").Collection("users")
		if err != nil {
			panic(err)
		}

		doc := models.RegisterUser{Username: user.Username, Password: user.Password}
		result, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			panic(err)
		}

		c.JSON(200, result)
	})

	e.POST("/login", func(c *gin.Context) {
		var user models.RegisterUser
		err := c.Bind(&user)
		if err != nil {
			panic(err)
		}

		if user.Username == "" || user.Password == "" {
			panic("User needs to enter bot ha username and a password")
		}

		connectionURI := os.Getenv("mongo_uri")
		serverApi := options.ServerAPI(options.ServerAPIVersion1)
		err = godotenv.Load()
		if err != nil {
			panic(err)
		}

		opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverApi)

		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			panic(err)
		}

		coll, err := client.Database("percent-back-app").Collection("users").Find(context.TODO(), bson.M{"username": user.Username})
		if err != nil {
			panic(err)
		}
		defer coll.Close(context.TODO())
		// for coll.Next(context.TODO()) {
		// 	var user
		// }

	})

}