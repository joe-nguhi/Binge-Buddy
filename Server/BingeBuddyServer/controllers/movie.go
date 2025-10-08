package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/database"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/models"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

const DefaultFavoriteMoviesLimit = 5

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var movies []models.Movie
		var movieCollection = database.OpenCollection("movies", client)
		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			log.Printf("Error fetching movies1: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching movies",
			})
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &movies); err != nil {
			log.Printf("Error fetching movies: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching movies",
			})
		}

		c.JSON(200, gin.H{
			"movies": movies,
		})
	}
}

func GetMovie(client *mongo.Client) gin.HandlerFunc {
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
		var movieCollection = database.OpenCollection("movies", client)
		result := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID})

		if err := result.Decode(&movie); err != nil {

			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Movie not found",
				})
				return
			}

			log.Printf("Error Decoding Movie: %v\n", err)

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

func AddMovie(client *mongo.Client) gin.HandlerFunc {
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
		var movieCollection = database.OpenCollection("movies", client)
		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			fmt.Printf("Error inserting movie: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Movie added successfully", "inserted_id": result.InsertedID})
	}

}

func UpdateAdminReview(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		role, err := utils.GetUserRoleFromContext(c)
		if err != nil {
			log.Printf("Error fetching user: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user role not defined"})
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}

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
		rankingName, rankingScore, err := getReviewRanking(req.AdminReview, client)

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
		var movieCollection = database.OpenCollection("movies", client)
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

func GetRankings(client *mongo.Client) ([]models.Ranking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var rankings []models.Ranking
	var rankingCollection = database.OpenCollection("rankings", client)
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

func getReviewRanking(review string, client *mongo.Client) (string, int, error) {
	godotenv.Load()

	rankings, err := GetRankings(client)

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

func GetMovieRecommendations(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)

		if err != nil {
			log.Printf("Error fetching user: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
			return
		}

		favoriteGenres := getUserFavoriteGenres(userId, client)

		var moviesLimit int64 = DefaultFavoriteMoviesLimit

		if err := godotenv.Load(".env"); err != nil {
			log.Println("Warning! unable to find .env file", err)
		}

		if limit := os.Getenv("FAVORITE_MOVIES_LIMIT"); limit != "" {
			moviesLimit, err = strconv.ParseInt(limit, 10, 64)
			if err != nil {
				log.Println("Warning! unable to parse FAVORITE_MOVIES_LIMIT", err)
			}
		}

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"ranking.ranking_value", 1}})
		findOptions.SetLimit(moviesLimit)

		filter := bson.M{"genre.genre_name": bson.M{"$in": favoriteGenres}}

		fmt.Printf("Filter: %v\n", filter)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var movieCollection = database.OpenCollection("movies", client)
		cursor, err := movieCollection.Find(ctx, filter, findOptions)

		if err != nil {
			log.Printf("Error fetching movies: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching movies",
			})
			return
		}

		defer cursor.Close(ctx)

		var movies []models.Movie

		if err := cursor.All(ctx, &movies); err != nil {
			log.Printf("Error decoding movies: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error fetching movies",
			})
			return
		}

		c.JSON(200, gin.H{
			"movies": movies,
		})
	}
}

func getUserFavoriteGenres(userId string, client *mongo.Client) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	var userCollection = database.OpenCollection("users", client)
	result := userCollection.FindOne(ctx, bson.M{"user_id": userId})
	if err := result.Decode(&user); err != nil {
		log.Printf("Error fetching user: %v\n", err)
		return []string{}
	}

	favs := make([]string, 0, len(user.FavoriteGenres))

	for _, genre := range user.FavoriteGenres {
		favs = append(favs, genre.GenreName)
	}

	return favs
}
