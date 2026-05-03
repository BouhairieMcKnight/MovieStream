package utils

import (
	"context"
	"crypto/rand"
	"errors"
	"os"
	"time"

	"github.com/BouhairieMcKnight/MovieStream/Server/persistence"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email:    email,
		LastName: lastName,
		Role:     role,
		UserId:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MovieStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken := rand.Text()

	return signedToken, refreshToken, nil
}

func UpdateRefreshToken(userId, refreshToken string, client *mongo.Client, c *gin.Context) (err error) {
	var ctx, cancel = context.WithTimeout(c, 10*time.Second)
	defer cancel()

	expiresAt, _ := time.Parse(time.RFC3339, time.Now().UTC().Add(24*30*time.Hour).Format(time.RFC3339))
	updateData := bson.M{
		"$push": bson.M{
			"refresh_tokens": bson.M{
				"token_id":  refreshToken,
				"expires_at": expiresAt,
			},
		},
	}

	var userCollection *mongo.Collection = persistence.OpenCollection("users", client)

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)

	if err != nil {
		return err
	}

	return nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}

	tokenString := authHeader[len("Bearer "):]

	if tokenString == "" {
		return "", errors.New("Bearer token is required")
	}

	return tokenString, nil
}

func GetCookieToken(token string, c *gin.Context) (string, error) {
	result, err := c.Cookie(token)

	if err != nil {
		return "", err
	}

	return result, nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	return claims, nil
}

func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")

	if !exists {
		return "", errors.New("userId does not exist in this context")
	}

	id, ok := userId.(string)
	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return id, nil
}

func GetUserRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")

	if !exists {
		return "", errors.New("userId does not exist in this context")
	}

	userRole, ok := role.(string)
	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return userRole, nil
}

func ValidateRefreshToken(tokenString string, userId string, client *mongo.Client, c *gin.Context) error {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var userCollection *mongo.Collection = persistence.OpenCollection("users", client)

	date := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.D{
    	{Key: "user_id", Value: userId},
    	{Key: "refresh_tokens", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "expires_at", Value: bson.D{
					{Key: "$gt", Value: date}}},
				{Key: "token_id", Value: tokenString}},
			}},
		},
	}

	projection := bson.D{{Key: "user_id", Value: userId}}
	opts := options.FindOne().SetProjection(projection)

	err := userCollection.FindOne(ctx, filter, opts).Err()
	if err != nil {
		return errors.New("Unable to find valid refresh tokens")
	}

	return nil
}

func InvalidateTokens(userId, tokenString string, client *mongo.Client, c *gin.Context) error {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}
	now, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
	update := bson.M{
		"$pull": bson.M{
			"refresh_tokens": bson.M{
				"$or": bson.A{
					bson.M{"expires_at": bson.M{"$lte": now}},
					bson.M{"token_id": tokenString},
				}},
			},
		}


	var userCollection = persistence.OpenCollection("users", client)
	result, err := userCollection.UpdateMany(ctx, filter, update)

	if err != nil || result.ModifiedCount == 0 {
		return err
	}

	return nil
}
