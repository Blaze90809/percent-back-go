package routes

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
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
	e.POST("/forgot-password", RateLimitMiddleware(SensitiveLimiter), func(c *gin.Context) {
		var body struct {
			Email string `json:"email"`
		}

		if err := c.BindJSON(&body); err != nil {
			fmt.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user models.User
		coll := client.Database("percent-back-app").Collection("users")
		err := coll.FindOne(context.TODO(), bson.M{"username": body.Email}).Decode(&user)
		if err != nil {
			fmt.Println("User not found or error finding user:", err)
			// Prevent user enumeration: always return a successful 200 status with generic message
			c.JSON(http.StatusOK, gin.H{"message": "If that email is registered, we have sent a reset link to it"})
			return
		}

		tokenBytes := make([]byte, 32)
		_, err = rand.Read(tokenBytes)
		if err != nil {
			fmt.Println("Failed to generate secure token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}
		resetToken := hex.EncodeToString(tokenBytes)

		// Hash the token using SHA-256 for secure database storage
		hasher := sha256.New()
		hasher.Write([]byte(resetToken))
		hashedToken := hex.EncodeToString(hasher.Sum(nil))

		update := bson.M{
			"$set": bson.M{
				"passwordResetToken":   hashedToken,
				"passwordResetExpires": time.Now().Add(1 * time.Hour),
			},
		}

		_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
		if err != nil {
			fmt.Println("Failed to update reset token in database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}

		resendClient := resend.NewClient(os.Getenv("RESEND_API_KEY"))

		resetURL := "https://nordicracetrack.com/reset-password/" + resetToken

		params := &resend.SendEmailRequest{
			From:    "blaze@nordicracetrack.com",
			To:      []string{user.Username},
			Subject: "Password Reset Request",
			Html:    "<strong>Click the link to reset your password:</strong> <a href=\"" + resetURL + "\">Reset Password</a>",
		}

		_, err = resendClient.Emails.Send(params)
		if err != nil {
			fmt.Println("Resend failed to send email:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "If that email is registered, we have sent a reset link to it"})
	})

	e.POST("/reset-password", RateLimitMiddleware(SensitiveLimiter), func(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token and password are required"})
			return
		}

		// Enforce password strength minimum length (8 characters)
		if len(body.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
			return
		}

		// Hash the incoming token using SHA-256 to compare with database
		hasher := sha256.New()
		hasher.Write([]byte(body.Token))
		hashedToken := hex.EncodeToString(hasher.Sum(nil))

		var user models.User
		coll := client.Database("percent-back-app").Collection("users")
		err := coll.FindOne(context.TODO(), bson.M{
			"passwordResetToken":   hashedToken,
			"passwordResetExpires": bson.M{"$gt": time.Now()},
		}).Decode(&user)
		if err != nil {
			fmt.Println("Error finding user with token:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Failed to hash password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
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
			fmt.Println("Failed to update user password in database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	})
}