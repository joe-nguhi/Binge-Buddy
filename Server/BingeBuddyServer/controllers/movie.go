package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/database"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var movieCollection = database.OpenCollection("movies")

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching movies1: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching movies",
			})
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &movies); err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching movies: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching movies",
			})
		}

		c.JSON(200, gin.H{
			"movies": movies,
		})
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		}

		var movie models.Movie

		result := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID})

		if err := result.Decode(&movie); err != nil {
			fmt.Fprintf(os.Stderr, "Error Decoding Movie: %v\n", err)

			c.JSON(http.StatusNotFound, gin.H{
				"error": "Movie not found",
			})

			return
		}

		c.JSON(200, gin.H{
			"movie": movie,
		})
	}
}
