package models

type RegisterUser struct {
	Username 	string 	`json:"Username" bson:"Username"`
	Password	string	`json:"Password" bson:"Password"`
}