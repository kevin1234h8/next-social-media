package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"social/project/database"
	"social/project/initializer"
	"social/project/model"
	"social/project/util"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func IsEmailExists(email string) bool {
	client, err := database.Connect()
	if err != nil {
		log.Fatal(err)
		return true
	}
	userCollections := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("USER_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := userCollections.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		log.Fatal(err)
		return true
	}

	return count > 0
}

func GetUserByEmail(email string) (model.User, error) {
	var user model.User
	client, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	userCollections := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("USER_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = userCollections.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Fatal("error find data from collections ", err)
	}
	if err != nil {
		return user, err
	}

	return user, nil

}

func Register(c *gin.Context) {
	body := model.User{}
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	hashedPassword, err := util.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	if IsEmailExists(body.Email) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
		return
	}

	user := model.User{
		Id:        primitive.NewObjectID(),
		Username:  body.Username,
		Email:     body.Email,
		Password:  hashedPassword,
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	logger := initializer.InitializeLogger()
	client, err := database.Connect()
	if err != nil {
		logger.Fatal(err.Error())
	}
	userCollections := client.Database("socialmedianext").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = userCollections.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create new user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := GetUserByEmail(credentials.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	token := createJWTToken(user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func GetUserInfo(c *gin.Context) {
	userClaims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User claims not found",
		})
		return
	}

	// Now you can use the userClaims to access user information
	userID := userClaims.(jwt.MapClaims)["user_id"].(string)
	email := userClaims.(jwt.MapClaims)["email"].(string)
	

	// Perform actions with userID
	c.JSON(http.StatusOK, gin.H{
		"message": "Protected route accessed",
		"userId":  userID,
		"email":   email,
	})
}

func createJWTToken(user model.User) string {
	claims := jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration (24 hours)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("SECRET_KEY")) // Use your own secret key

	return tokenString
}
