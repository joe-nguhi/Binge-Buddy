package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/database"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/models"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data", "details": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()
		var userCollection = database.OpenCollection("users", client)
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong. Try again later"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "An Account by this Email already exists!"})
			return
		}

		user.Password, err = hashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong. Try again later"})
			return
		}

		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		result, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong. Try again later"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User Registered Successfully", "id": result.InsertedID})
	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func LoginUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var formData models.UserLogin

		if err := c.ShouldBindJSON(&formData); err != nil {
			fmt.Printf("Invalid Data: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
			return
		}

		if err := validate.Struct(formData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data", "details": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()
		var userCollection = database.OpenCollection("users", client)
		found := userCollection.FindOne(ctx, bson.M{"email": formData.Email})

		var user models.User

		if err := found.Decode(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email or Password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formData.Password)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email or Password"})
			return
		}

		authToken, refreshToken, err := utils.GenerateUserTokens(user.UserID, userCollection, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong. Try again later"})
			return
		}

		if err := utils.UpdateUserTokens(user.UserID, authToken, refreshToken, userCollection, c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong. Try again later"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login Successful", "user": models.UserResponse{
			UserID:         user.UserID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Email:          user.Email,
			Role:           user.Role,
			FavoriteGenres: user.FavoriteGenres,
			Token:          authToken,
			RefreshToken:   refreshToken,
		}})

	}
}
