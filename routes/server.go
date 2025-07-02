package routes

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRouter() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	connectionURI := os.Getenv("mongo_uri")
	opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	e := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	e.Use(cors.New(config))

	e.GET("/", func(c *gin.Context) {
		c.String(200, "Nordic percent back.")
	})

	racesRoutes(e, client)
	usersRoutes(e, client)
	loginRoutes(e, client)
	passwordResetRoutes(e, client)

	err = e.Run()
	if err != nil {
		panic(err)
	}
}
