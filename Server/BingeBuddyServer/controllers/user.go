package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/database"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Move client init to main and pass it to controllers
var userCollection *mongo.Collection = database.OpenCollection("users")

func RegisterUser() gin.HandlerFunc {
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

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

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
