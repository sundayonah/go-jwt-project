package helpers

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sundayonah/go-jwt-project/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// type SignedDetaills struct {
// 	Email     string
// 	FirstName string
// 	LastName  string
// 	UserId       string
// 	UserType  string
// 	jwt.StandardClaims
// }

// var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// var SECRET_KEY string = os.Getenv("SECRET_KEY")

// func GenerateAllTokens(email string, firstName string, lastName string, userType string, UserId string) (signedToken string, signedRefreshToken string, err error) {
// 	claims := &SignedDetaills{
// 		Email:     email,
// 		FirstName: firstName,
// 		LastName:  lastName,
// 		UserId:       UserId,
// 		UserType:  userType,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
// 		},
// 	}

// 	refreshClaims := &SignedDetaills{
// 		ExpiresAt: time.Now().Add(time.Hour * time.Duration(168)).Unix()}

// 	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

// 	fmt.Println("Access Token:", token)
// 	fmt.Println("Refresh Token:", refreshToken)

// 	if err != nil{
// 		log.Panic(err)
// 		return
// 	}

// 	return token, refreshToken, err
// }

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserId    string
	UserType  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, userType, userId string) (signedToken string, signedRefreshToken string, err error) {
	// Creating the main token claims
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserId:    userId,
		UserType:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	// Creating the refresh token claims
	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserId:    userId,
		UserType:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(168 * time.Hour).Unix(), // 7 days
		},
	}

	// Generating the main token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic("Error generating token: ", err)
		return "", "", err
	}

	// Generating the refresh token
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic("Error generating refresh token: ", err)
		return "", "", err
	}

	// Return both tokens
	return token, refreshToken, nil
}

func UpdateAllTokens(signedToken, signedRefreshToken, userId string) error {
	// Set up a context with a timeout of 100 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Create the update object with the new token, refresh token, and updated_at timestamp
	updateObj := bson.D{
		{Key: "token", Value: signedToken},
		{Key: "refresh_token", Value: signedRefreshToken},
		{Key: "updated_at", Value: time.Now()},
	}

	// Specify the filter and upsert option
	filter := bson.M{"user_id": userId}
	upsert := true
	opt := options.Update().SetUpsert(upsert)

	// Perform the update operation
	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, opt)
	if err != nil {
		return err // Return the error to allow the caller to handle it
	}

	return nil // Return nil if the operation is successful
}

// func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	var updateObj primitive.D

// 	updateObj = append(updateObj, bson.E{"token", signedToken})
// 	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

// 	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

// 	upsert := true
// 	filter := bson.M{"user_id": userId}
// 	opt := options.UpdateOptions{
// 		Upsert: &upsert,
// 	}
// 	_, err := userCollection.UpdateOne(ctx, filter, bson.D{"$set": updateObj}, &opt)

// 	if err != nil {
// 		log.Panic(err)
// 		return
// 	}

// 	return

// }
