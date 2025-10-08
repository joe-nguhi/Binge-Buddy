package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joe-nguhi/Binge-Buddy/Server/BingeBuddyServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type signingData struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
	Role      string
}

type SignedData struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
	Role      string
	jwt.RegisteredClaims
}

func GenerateUserTokens(userID string, userCollection *mongo.Collection, c context.Context) (string, string, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	found := userCollection.FindOne(ctx, bson.M{"user_id": userID})

	var user models.User

	if err := found.Decode(&user); err != nil {
		return "", "", err
	}

	signData := signingData{
		user.UserID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
	}

	authToken, err := generateToken(signData, "JWT_KEY", time.Hour*24)
	refreshToken, err := generateToken(signData, "JWT_REFRESH_KEY", time.Hour*24*7)

	if err != nil {
		return "", "", err
	}

	return authToken, refreshToken, err
}

func generateToken(data signingData, keyEnvVar string, expiration time.Duration) (string, error) {
	signingKey := []byte(os.Getenv(keyEnvVar))
	claim := SignedData{
		data.UserID,
		data.Email,
		data.FirstName,
		data.LastName,
		data.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Binge Buddy",
			Subject:   data.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	ss, err := token.SignedString(signingKey)

	return ss, err
}

func UpdateUserTokens(userID, authToken, refreshToken string, userCollection *mongo.Collection, c context.Context) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"token": authToken, "refresh_token": refreshToken, "updated_at": time.Now()}})

	return err
}

func ValidateToken(signedToken string) (claims *SignedData, err error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	}
	claims = &SignedData{}

	token, err := jwt.ParseWithClaims(signedToken, claims, keyFunc)

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("signing method invalid")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil

}
