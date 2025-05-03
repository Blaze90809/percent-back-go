package models


import ("go.mongodb.org/mongo-driver/bson/primitive"
"github.com/golang-jwt/jwt/v5")

type RegisterUser struct {
	Username 	string 	`json:"Username" bson:"Username"`
	Password	string	`json:"Password" bson:"Password"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	V        int                `bson:"__v"`
}

type Claims struct {
	UserID  primitive.ObjectID  `json:"userId"`
	jwt.RegisteredClaims
}