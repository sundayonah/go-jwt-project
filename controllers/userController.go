package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sundayonah/go-jwt-project/database"
	"github.com/sundayonah/go-jwt-project/helpers"
	"github.com/sundayonah/go-jwt-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		return false, "Invalid Email or Password"
	}
	return true, ""

}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// Bind JSON to user model
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Hash password before saving to the database
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while hashing password"})
			return
		}
		user.Password = hashedPassword

		// Check if email or phone already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{
			"$or": []bson.M{
				{"email": user.Email},
				{"phone": user.Phone},
			},
		})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for email or phone"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone already exists"})
			return
		}

		// Set timestamps and IDs
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// Generate tokens
		token, refreshToken, tokenErr := helpers.GenerateAllTokens(user.Email, user.FirstName, user.LastName, *&user.User_type, user.User_id)
		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate tokens"})
			return
		}
		user.Token = &token
		user.Refresh_token = &refreshToken

		// Insert user into the database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "user created successfully", "user_id": resultInsertionNumber.InsertedID})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		isPasswordValid, msg := VerifyPassword(user.Password, foundUser.Password)
		if !isPasswordValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		// Generate tokens if authentication is successful
		token, refreshToken, tokenErr := helpers.GenerateAllTokens(
			foundUser.Email,
			foundUser.FirstName,
			foundUser.LastName,
			foundUser.User_type,
			foundUser.User_id,
		)

		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		// Use UpdateAllTokens function to update tokens in the database
		if err := helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user tokens"})
			return
		}

		// Send the generated tokens in the response
		c.JSON(http.StatusOK, gin.H{
			"message":       "Login successful",
			"token":         token,
			"refresh_token": refreshToken,
		})


		// // Update user tokens in the database
		// updateFields := bson.M{
		// 	"token":         token,
		// 	"refresh_token": refreshToken,
		// 	"updated_at":    time.Now(),
		// }

		// _, err = userCollection.UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": updateFields})
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user tokens"})
		// 	return
		// }

		// // Send the generated tokens in the response
		// c.JSON(http.StatusOK, gin.H{
		// 	"message":       "Login successful",
		// 	"token":         token,
		// 	"refresh_token": refreshToken,
		// })

	}

}

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)

	}
}

// func GetUsersByRole()
