package routes

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setNoCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

// SetupRouter builds and returns the configured Gin engine.
func SetupRouter(client *mongo.Client, staticDir string) *gin.Engine {
	if staticDir == "" {
		staticDir = "./static"
	}

	e := gin.Default()

	// Security Headers Middleware
	e.Use(func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://*.ondigitalocean.app"} // Allow local and DO domains
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	e.Use(cors.New(config))

	// The React build puts assets in a nested 'static' folder.
	// We point the /static route to ./static/static to match the browser's requests.
	e.Static("/static", filepath.Join(staticDir, "static"))

	// Middleware to prevent caching of index.html
	noCache := func(c *gin.Context) {
		setNoCacheHeaders(c)
		c.Next()
	}

	indexFile := filepath.Join(staticDir, "index.html")

	e.GET("/", noCache, func(c *gin.Context) {
		c.File(indexFile)
	})
	e.StaticFile("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))
	e.StaticFile("/manifest.json", filepath.Join(staticDir, "manifest.json"))
	e.StaticFile("/robots.txt", filepath.Join(staticDir, "robots.txt"))
	e.StaticFile("/logo192.png", filepath.Join(staticDir, "logo192.png"))
	e.StaticFile("/logo512.png", filepath.Join(staticDir, "logo512.png"))
	e.StaticFile("/20230313_172031.jpg", filepath.Join(staticDir, "20230313_172031.jpg"))

	api := e.Group("/api")
	api.Use(RateLimitMiddleware(GeneralLimiter))
	{
		racesRoutes(api, client)
		usersRoutes(api, client)
		loginRoutes(api, client)
		passwordResetRoutes(api, client)
	}

	e.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		if strings.HasPrefix(reqPath, "/api/") || reqPath == "/api" ||
			strings.Contains(filepath.Base(reqPath), ".") || strings.Contains(reqPath, "/.") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		setNoCacheHeaders(c)
		c.File(indexFile)
	})

	return e
}

func NewRouter() {
	// Only load .env in development. In production (DigitalOcean), environment variables are set directly.
	_ = godotenv.Load()

	// Validate environment variables
	if os.Getenv("JWT_SECRET_KEY") == "" {
		panic("JWT_SECRET_KEY is not set")
	}
	connectionURI := os.Getenv("mongo_uri")
	if connectionURI == "" {
		panic("mongo_uri is not set")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	e := SetupRouter(client, "./static")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = e.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
