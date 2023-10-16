package server

import "github.com/gin-gonic/gin"

func NewRouter() {
	e := gin.Default()

	e.GET("/", func(c *gin.Context) {
		c.String(200, "Nordic percent back.")
	})

	err := e.Run()
	if err != nil {
		panic(err)
	}
}