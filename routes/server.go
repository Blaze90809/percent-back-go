package routes

import (
	"context"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRouter() {
	// Only load .env in development. In production (DigitalOcean), environment variables are set directly.
	_ = godotenv.Load()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	connectionURI := os.Getenv("mongo_uri")
	opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	e := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://*.ondigitalocean.app"} // Allow local and DO domains
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	e.Use(cors.New(config))

	// The React build puts assets in a nested 'static' folder. 
	// We point the /static route to ./static/static to match the browser's requests.
	e.Static("/static", "./static/static")

	// Middleware to prevent caching of index.html
	noCache := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}

	e.GET("/", noCache, func(c *gin.Context) {
		c.File("./static/index.html")
	})
	e.StaticFile("/favicon.ico", "./static/favicon.ico")
	e.StaticFile("/manifest.json", "./static/manifest.json")
	e.StaticFile("/robots.txt", "./static/robots.txt")
	e.StaticFile("/logo192.png", "./static/logo192.png")
	e.StaticFile("/logo512.png", "./static/logo512.png")
	e.StaticFile("/20230313_172031.jpg", "./static/20230313_172031.jpg")

	api := e.Group("/api")
	{
		racesRoutes(api, client)
		usersRoutes(api, client)
		loginRoutes(api, client)
		passwordResetRoutes(api, client)
	}

	e.NoRoute(noCache, func(c *gin.Context) {
		c.File("./static/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = e.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
