package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultMaxBodyBytes int64 = 1 << 20 // 1 MiB

var defaultTrustedProxies = []string{
	"127.0.0.1/32",
	"::1/128",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

func setNoCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

func parseCSVEnv(raw string) []string {
	var result []string
	for _, item := range strings.Split(raw, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func configureTrustedProxies(e *gin.Engine) {
	if platform := strings.TrimSpace(os.Getenv("TRUSTED_PLATFORM")); platform != "" {
		e.TrustedPlatform = platform
	}

	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	var proxies []string
	switch {
	case strings.EqualFold(raw, "none"):
		proxies = nil
	case raw != "":
		proxies = parseCSVEnv(raw)
	default:
		proxies = defaultTrustedProxies
	}

	if err := e.SetTrustedProxies(proxies); err != nil {
		panic(fmt.Sprintf("invalid TRUSTED_PROXIES configuration: %v", err))
	}
}

func allowedCORSOrigins() []string {
	if origins := parseCSVEnv(os.Getenv("ALLOWED_ORIGINS")); len(origins) > 0 {
		return origins
	}
	if gin.Mode() != gin.ReleaseMode {
		return []string{"http://localhost:3000"}
	}
	return nil
}

// MaxBodySizeMiddleware caps request body size to protect against memory exhaustion.
func MaxBodySizeMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body too large",
			})
			return
		}
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}

// NewHTTPServer creates an http.Server with defensive timeouts to mitigate slow-client DoS.
func NewHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

// SetupRouter builds and returns the configured Gin engine.
func SetupRouter(client *mongo.Client, staticDir string) *gin.Engine {
	if staticDir == "" {
		staticDir = "./static"
	}

	e := gin.Default()
	configureTrustedProxies(e)

	// Security Headers Middleware
	e.Use(func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})

	e.Use(MaxBodySizeMiddleware(defaultMaxBodyBytes))

	if origins := allowedCORSOrigins(); len(origins) > 0 {
		config := cors.DefaultConfig()
		config.AllowOrigins = origins
		config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
		e.Use(cors.New(config))
	}

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

	srv := NewHTTPServer(":"+port, e)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
