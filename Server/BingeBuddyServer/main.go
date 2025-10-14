package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
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

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	//config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	var client = database.Connection()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}(client, context.Background())

	routes.SetupRoutes(router, client)

	err := router.Run(port)

	if err != nil {
		fmt.Println("Failed to start Server", err)
	}
}
