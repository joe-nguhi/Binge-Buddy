package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	controller "github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/controllers"
)

const port = ":8010"

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Binge Buddy, 👏")
	})

	r.GET("/movies", controller.GetMovies())
	r.GET("/movie/:imdb_id", controller.GetMovie())
	r.POST("/movie/add", controller.AddMovie())

	err := r.Run(port)

	if err != nil {
		fmt.Println("Failed to start Server", err)
	}
}
