package routes

import (
	controller "github.com/BouhairieMcKnight/MovieStream/Server/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


func SetupUnprotectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.GET("/movies/movies", controller.GetMovies(client))
	router.POST("/users/register", controller.RegisterUser(client))
	router.GET("/movies/search", controller.SearchMovies(client))
	router.GET("/movies/trending", controller.GetTrendingMovies(client))
	router.POST("/users/login", controller.LoginUser(client))
	router.POST("/users/logout", controller.Logout(client))
	router.GET("/movies/genres", controller.GetGenres(client))
	router.POST("/users/refresh", controller.RefreshTokenHandler(client))
}