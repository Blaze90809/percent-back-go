package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Race struct {
	ID           primitive.ObjectID `bson:"_id"`
	RaceName     string             `bson:"raceName"`
	RaceDate     string             `bson:"raceDate"`
	RaceDistance float64            `bson:"raceDistance"`
	PercentBack  float64            `bson:"percentBack"`
}

type CreateRace struct {
	RaceName     string             `json:"raceName" bson:"raceName"`
	RaceDate     string             `json:"raceDate" bson:"raceDate"`
	RaceDistance float64            `json:"raceDistance" bson:"raceDistance"`
	PercentBack  float64            `json:"percentBack" bson:"percentBack"`
	UserID       primitive.ObjectID `json:"userId" bson:"userId"`
}