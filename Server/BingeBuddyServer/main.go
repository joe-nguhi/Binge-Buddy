package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/database"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const port = ":8010"

func main() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Binge Buddy, 👏")
	})

	var client *mongo.Client = database.Connection()

	routes.SetupRoutes(router, client)

	err := router.Run(port)

	if err != nil {
		fmt.Println("Failed to start Server", err)
	}
}
