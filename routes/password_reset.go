package routes

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"react-app-golang/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func passwordResetRoutes(e *gin.RouterGroup, client *mongo.Client) {
	e.POST("/forgot-password", func(c *gin.Context) {
		var body struct {
			Email string `json:"email"`
		}

		if err := c.BindJSON(&body); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user models.User
		coll := client.Database("percent-back-app").Collection("users")
		err := coll.FindOne(context.TODO(), bson.M{"username": body.Email}).Decode(&user)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		tokenBytes := make([]byte, 32)
		_, err = rand.Read(tokenBytes)
		if err != nil {
						fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		resetToken := hex.EncodeToString(tokenBytes)

		update := bson.M{
			"$set": bson.M{
				"passwordResetToken":   resetToken,
				"passwordResetExpires": time.Now().Add(1 * time.Hour),
			},
		}

		_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
			return
		}

		client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

		resetURL := "https://nordicracetrack.com/reset-password/" + resetToken

		params := &resend.SendEmailRequest{
			From:    "blaze@nordicracetrack.com",
			To:      []string{user.Username},
			Subject: "Password Reset Request",
			Html:    "<strong>Click the link to reset your password:</strong> <a href=\"" + resetURL + "\">Reset Password</a>",
		}

		_, err = client.Emails.Send(params)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
	})

	e.POST("/reset-password", func(c *gin.Context) {
		var body struct {
			Token    string `json:"token"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&body); err != nil {
			fmt.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if body.Token == "" || body.Password == "" {
			fmt.Println("Token or Password missing in request body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token and password are required"})
			return
		}

		var user models.User
		coll := client.Database("percent-back-app").Collection("users")
		err := coll.FindOne(context.TODO(), bson.M{"passwordResetToken": body.Token, "passwordResetExpires": bson.M{"$gt": time.Now()}}).Decode(&user)
		if err != nil {
			fmt.Println("Error finding user with token:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"password": string(hashedPassword),
			},
			"$unset": bson.M{
				"passwordResetToken":   "",
				"passwordResetExpires": "",
			},
		}

		_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	})
}