package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() {
	e := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	e.Use(cors.New(config))

	e.GET("/", func(c *gin.Context) {
		c.String(200, "Nordic percent back.")
	})

	racesRoutes(e)
	usersRoutes(e)
	loginRoutes(e)

	err := e.Run()
	if err != nil {
		panic(err)
	}
}
