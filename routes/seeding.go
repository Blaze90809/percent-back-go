package routes

import (
	"context"
	"net/http"

	"react-app-golang/internal/seeding"
	"react-app-golang/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EvaluateRequest struct {
	RaceName     string  `json:"raceName"`
	RaceDate     string  `json:"raceDate"`
	RaceDistance float64 `json:"raceDistance"`
	PercentBack  float64 `json:"percentBack"`
	QualifierID  string  `json:"qualifierId"`
}

func seedingRoutes(e *gin.RouterGroup, client *mongo.Client) {
	sg := e.Group("/seeding")

	// GET /api/seeding/qualifiers - returns static qualifiers list
	sg.GET("/qualifiers", func(c *gin.Context) {
		store := seeding.GetStore()
		c.JSON(http.StatusOK, store.Qualifiers)
	})

	// GET /api/seeding/cutoffs - returns static cutoffs by event
	sg.GET("/cutoffs", func(c *gin.Context) {
		store := seeding.GetStore()
		c.JSON(http.StatusOK, store.Cutoffs)
	})

	// POST /api/seeding/evaluate - evaluates ad-hoc race/percent back across all 4 events
	sg.POST("/evaluate", func(c *gin.Context) {
		var req EvaluateRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		result := seeding.EvaluateRaceSeeding(seeding.RaceInputData{
			RaceName:     req.RaceName,
			RaceDate:     req.RaceDate,
			RaceDistance: req.RaceDistance,
			PercentBack:  req.PercentBack,
		}, req.QualifierID)

		c.JSON(http.StatusOK, result)
	})

	// GET /api/seeding/my-waves - evaluates stored races by season for authenticated user
	sg.GET("/my-waves", authMiddleware(), func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		coll, err := client.Database("percent-back-app").Collection("races").Find(context.TODO(), bson.M{"userId": userId})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer coll.Close(context.TODO())

		races := make([]models.Race, 0)
		for coll.Next(context.TODO()) {
			var race models.Race
			if err := coll.Decode(&race); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			races = append(races, race)
		}

		inputRaces := make([]seeding.RaceInputData, 0, len(races))
		for _, r := range races {
			inputRaces = append(inputRaces, seeding.RaceInputData{
				ID:           r.ID.Hex(),
				RaceName:     r.RaceName,
				RaceDate:     r.RaceDate,
				RaceDistance: r.RaceDistance,
				PercentBack:  r.PercentBack,
			})
		}

		summaries := seeding.GroupAndEvaluateSeasons(inputRaces)
		c.JSON(http.StatusOK, summaries)
	})
}
