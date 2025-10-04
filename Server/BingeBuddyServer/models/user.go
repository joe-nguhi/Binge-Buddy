package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID         string        `bson:"user_id" json:"user_id" `
	FirstName      string        `bson:"first_name" json:"first_name" validate:"required"`
	LastName       string        `bson:"last_name" json:"last_name" validate:"required"`
	Email          string        `bson:"email" json:"email" validate:"required,email"`
	Password       string        `bson:"password" json:"password" validate:"required"`
	Role           string        `bson:"role" json:"role" validate:"oneof=ADMIN USER"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
	Token          string        `bson:"token" json:"token"`
	RefreshToken   string        `bson:"refresh_token" json:"refresh_token"`
	FavoriteGenres []Genre       `bson:"favorite_genres" json:"favorite_genres"`
}
