package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BouhairieMcKnight/MovieStream/Server/domain"
	"github.com/BouhairieMcKnight/MovieStream/Server/dtos"
	"github.com/BouhairieMcKnight/MovieStream/Server/persistence"
	"github.com/BouhairieMcKnight/MovieStream/Server/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var validate = validator.New()

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = persistence.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, bson.D{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		defer cursor.Close(ctx)

		var movies []domain.Movie

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()


		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movieCollection *mongo.Collection = persistence.OpenCollection("movies", client)

		var movie domain.Movie

		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movie domain.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		var movieCollection *mongo.Collection = persistence.OpenCollection("movies", client)

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func AdminReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, err := utils.GetUserRoleFromContext(c); err != nil || role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Role"})
			return
		}

		movieId := c.Param("imdb_id")

		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id required"})
			return
		}

		var request dtos.AdminReview

		var response dtos.AdminReviewResponse

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		sentiment, rankVal, err := GetReviewRanking(request.AdminReview, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking"})
			return
		}

		filter := bson.M{"imdb_id": movieId}

		update := bson.M{
			"$set": bson.M{
				"admin_review": request.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = persistence.OpenCollection("movies", client)

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil || result.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}

		response.RankingName = sentiment
		response.AdminReview = request.AdminReview

		c.JSON(http.StatusOK, response)
	}
}

func GetReviewRanking(admin_review string, client *mongo.Client, c *gin.Context) (string, int, error) {
	rankings, err := GetRankings(client, c)

	if err != nil {
		return "", 0, err
	}

	sentimentDelimeted := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimeted = sentimentDelimeted + ranking.RankingName + ","
		}
	}

	sentimentDelimeted = strings.Trim(sentimentDelimeted, ",")

	err = godotenv.Load(".env")

	if err != nil {
		log.Println("Could not get environment variables")
	}

	openAiApiKey := os.Getenv("OPENAI_API_KEY")

	if openAiApiKey == "" {
		return "", 0, errors.New("Could not read OPENAI_API_KEY")
	}

	llm, err := openai.New(openai.WithToken(openAiApiKey))

	if err != nil {
		return "", 0, err
	}

	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimeted, 1)

	response, err := llm.Call(c, base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}

	rankVal := 0

	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil
}

func GetRankings(client *mongo.Client, c *gin.Context) ([]domain.Ranking, error) {
	var rankings []domain.Ranking

	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var rankingCollection *mongo.Collection = persistence.OpenCollection("rankings", client)

	cursor, err := rankingCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		favoriteGenres, err := GetUsersFavoriteGenres(userId, client, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		err = godotenv.Load(".env")

		if err != nil {
			log.Println("Warning: .env file not found")
		}

		var recommendedMovieLimitVal int64 = 5

		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")

		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}

		findOptions := options.Find()

		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})

		findOptions.SetLimit(recommendedMovieLimitVal)

		filter := bson.M{"genre.genre_name": bson.M{"$in": favoriteGenres}}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = persistence.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
		}
		defer cursor.Close(ctx)

		var recommendedMovies []domain.Movie

		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUsersFavoriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	projection := bson.M{
		"favorite_genres.genres_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)
	var result bson.M

	var userCollection *mongo.Collection = persistence.OpenCollection("users", client)

	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}

	favGenres, ok := result["favorite_genres"].(bson.A)

	if !ok {
		return []string{}, errors.New("unable to retrieve favorite geners for user")
	}

	var genreNames []string

	for _, item := range favGenres {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil
}

func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var genreCollection *mongo.Collection = persistence.OpenCollection("genres", client)

		cursor, err := genreCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie genres"})
			return
		}
		defer cursor.Close(ctx)

		var genres []domain.Genre
		if err := cursor.All(ctx, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, genres)

	}
}

func GetTrendingMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = persistence.OpenCollection("movies", client)

		filter := bson.D{}
		
		opts := options.Find().
			SetSort(bson.D{{Key: "ranking.ranking_value", Value: -1}}).
			SetLimit(5)

		cursor, err := movieCollection.Find(ctx, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching database", "details": err.Error()})
		}
		defer cursor.Close(ctx)

		var movies []domain.Movie
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not find any movies in database", "details": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, movies)
	}
}

func SearchMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c * gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieSize int64 = 15

		filter := bson.D{}
		opts := options.Find().SetLimit(movieSize)


		if sort := c.Query("page"); sort != "" {
			if sort == "asc" {
				opts.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
			}
			if sort == "desc" {
				opts.SetSort(bson.D{{Key: "ranking.ranking_value", Value: -1}})
			}
		}

		if term := c.Query("term"); term != "" {
			filter = append(filter, bson.E{Key: "title", Value: bson.D{
				{Key: "$regex", Value: term},
				{Key: "$options", Value: "i"},
			}})
		}
		if pageNum, err := strconv.ParseInt(c.Query("page"), 10, 64); err != nil {
			opts.SetSkip(pageNum * movieSize)
		}
		if genre := c.Query("genre"); genre != "" {
			filter = append(filter, bson.E{Key: "genre", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "genre_name", Value: genre},
				}},
			}})
		}

		var movieCollection *mongo.Collection = persistence.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, filter, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open database", "details": err.Error()})
		}

		defer cursor.Close(ctx)

		var movies []domain.Movie
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not find any movies in database", "details": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, movies)
	}
}

