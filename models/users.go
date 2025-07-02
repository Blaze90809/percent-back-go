package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUser struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id"`
	Username          string             `bson:"username"`
	Password          string             `bson:"password"`
	PasswordResetToken string            `bson:"passwordResetToken,omitempty"`
	PasswordResetExpires time.Time        `bson:"passwordResetExpires,omitempty"`
	V                 int                `bson:"__v"`
}

type Claims struct {
	UserID primitive.ObjectID `json:"userId"`
	jwt.RegisteredClaims
}