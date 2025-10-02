package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Binge Buddy, 👏")
	})

	err := r.Run(port)

	if err != nil {
		fmt.Println("Failed to start Server", err)
	}
}
