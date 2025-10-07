package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/controllers"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/middleware"
)

func SetupRoutes(app *gin.Engine) {

	app.GET("/movies", controllers.GetMovies())
	app.POST("/register", controllers.RegisterUser())
	app.POST("/login", controllers.LoginUser())

	app.Use(middleware.Authenticate())

	app.GET("/movie/:imdb_id", controllers.GetMovie())
	app.POST("/movie/add", controllers.AddMovie())
	app.PATCH("/movie/update-review/:imdb_id", controllers.UpdateAdminReview())

}
