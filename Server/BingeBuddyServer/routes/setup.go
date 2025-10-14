package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/controllers"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupRoutes(app *gin.Engine, client *mongo.Client) {

	app.GET("/movies", controllers.GetMovies(client))
	app.GET("/movies/genres", controllers.GetGenres(client))
	app.POST("/register", controllers.RegisterUser(client))
	app.POST("/login", controllers.LoginUser(client))

	app.Use(middleware.Authenticate())

	app.GET("/movies/recommended", controllers.GetMovieRecommendations(client))
	app.GET("/movie/:imdb_id", controllers.GetMovie(client))
	app.POST("/movie/add", controllers.AddMovie(client))
	app.PATCH("/movie/update-review/:imdb_id", controllers.UpdateAdminReview(client))

}
