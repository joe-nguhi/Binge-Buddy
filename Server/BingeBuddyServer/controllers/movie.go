package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/database"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection = database.OpenCollection("movies")
var rankingCollection = database.OpenCollection("rankings")
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

func UpdateAdminReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieId, _ := c.Params.Get("imdb_id")

		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Movie ID"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review" validate:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if err := validate.Struct(req); err != nil {
			fmt.Printf("Invalid Data: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data", "details": err.Error()})
			return
		}

		// Simulate Get Ranking from AI
		rankingName, rankingScore, err := GetReviewRanking(req.AdminReview)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ranking"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_name":  rankingName,
					"ranking_value": rankingScore,
				},
			},
		}

		result, err := movieCollection.UpdateOne(ctx, bson.M{"imdb_id": movieId}, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		var res struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		res.RankingName = rankingName
		res.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, res)
	}
}

func GetRankings() ([]models.Ranking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var rankings []models.Ranking

	cursor, err := rankingCollection.Find(ctx, bson.M{})

	if err != nil {
		return rankings, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rankings); err != nil {
		return rankings, err
	}

	return rankings, nil
}

func GetReviewRanking(review string) (string, int, error) {
	godotenv.Load()

	rankings, err := GetRankings()

	if err != nil {
		return "", 0, err
	}

	prompt := os.Getenv("BASE_PROMPT_TEMPLATE")

	sentiments := ""

	for _, ranking := range rankings {
		sentiments += fmt.Sprintf("%s, ", ranking.RankingName)
	}

	strings.Trim(sentiments, ", ")

	prompt = strings.Replace(prompt, "{{rankings}}", sentiments, -1)
	prompt = strings.Replace(prompt, "{{review}}", review, -1)

	fmt.Printf("Prompt: %s\n", prompt)

	// TODO: Call an AI API
	return "Excellent", 1, err
}
