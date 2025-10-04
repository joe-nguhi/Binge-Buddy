package utils

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/database"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var userCollection = database.OpenCollection("users")

type signingData struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
	Role      string
}

type signedData struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
	Role      string
	jwt.RegisteredClaims
}

func GenerateUserTokens(userID string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	found := userCollection.FindOne(ctx, bson.M{"user_id": userID})

	var user models.User

	if err := found.Decode(&user); err != nil {
		return "", "", err
	}

	userData := signingData{
		user.UserID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
	}

	authToken, err := generateAuthToken(userData)
	refreshToken, err := generateRefreshToken(userData)

	if err != nil {
		return "", "", err
	}

	return authToken, refreshToken, err
}

func generateAuthToken(data signingData) (string, error) {
	signingKey := []byte(os.Getenv("JWT_KEY"))
	claim := signedData{
		data.UserID,
		data.Email,
		data.FirstName,
		data.LastName,
		data.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Binge Buddy",
			Subject:   data.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	ss, err := token.SignedString(signingKey)

	return ss, err
}

func generateRefreshToken(data signingData) (string, error) {
	signingKey := []byte(os.Getenv("JWT_REFRESH_KEY"))
	claim := signedData{
		data.UserID,
		data.Email,
		data.FirstName,
		data.LastName,
		data.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Binge Buddy",
			Subject:   data.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	ss, err := token.SignedString(signingKey)

	return ss, err
}

func UpdateUserTokens(userID, authToken, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"token": authToken, "refresh_token": refreshToken, "updated_at": time.Now()}})

	return err
}
