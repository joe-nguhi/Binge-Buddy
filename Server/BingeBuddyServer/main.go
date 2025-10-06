package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/routes"
)

const port = ":8010"

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Binge Buddy, 👏")
	})

	routes.SetupRoutes(r)

	err := r.Run(port)

	if err != nil {
		fmt.Println("Failed to start Server", err)
	}
}
