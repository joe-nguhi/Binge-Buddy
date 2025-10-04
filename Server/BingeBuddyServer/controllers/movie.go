package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/database"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection = database.OpenCollection("movies")
var validate = validator.New(validator.WithRequiredStructEnabled())

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

			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Movie not found",
				})
				return
			}

			fmt.Fprintf(os.Stderr, "Error Decoding Movie: %v\n", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching movie",
			})

			return
		}

		c.JSON(200, gin.H{
			"movie": movie,
		})
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movie models.Movie

		if err := c.ShouldBindJSON(&movie); err != nil {
			fmt.Printf("Invalid Data: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			fmt.Printf("Invalid Data: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data", "details": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			fmt.Printf("Error inserting movie: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Movie added successfully", "inserted_id": result.InsertedID})
	}

}
