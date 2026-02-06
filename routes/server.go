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

	e.Static("/static", "./static")
	e.StaticFile("/", "./static/index.html")
	e.StaticFile("/favicon.ico", "./static/favicon.ico")
	e.StaticFile("/manifest.json", "./static/manifest.json")
	e.StaticFile("/robots.txt", "./static/robots.txt")
	e.StaticFile("/logo192.png", "./static/logo192.png")
	e.StaticFile("/logo512.png", "./static/logo512.png")

	api := e.Group("/api")
	{
		racesRoutes(api, client)
		usersRoutes(api, client)
		loginRoutes(api, client)
		passwordResetRoutes(api, client)
	}

	e.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	err = e.Run()
	if err != nil {
		panic(err)
	}
}
