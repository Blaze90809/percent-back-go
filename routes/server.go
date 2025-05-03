package routes

import (
	"github.com/gin-gonic/gin"
)



func NewRouter() {
	e := gin.Default()

	e.GET("/", func(c *gin.Context) {
		c.String(200, "Nordic percent back.")
	})

	racesRoutes(e)
	usersRoutes(e)

	err := e.Run()
	if err != nil {
		panic(err)
	}
}
