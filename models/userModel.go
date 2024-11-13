package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type User struct {
// 	ID            primitive.ObjectID `bson:"_id"`
// 	FirstName     string             `json:"first_name" validate:"required, min=2, max=100"`
// 	LastName      string             `json:"last_name" validate:"required, min=2, max=100"`
// 	Password      string             `json:"password" validate:"required, min=6"`
// 	Email         string             `json:"email" validate:"email, required"`
// 	Phone         string             `json:"phone" validate:"required"`
// 	Token         string             `json:"token"`
// 	User_type     string             `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
// 	Refresh_token string             `json:"refresh_token"`
// 	Created_at    time.Time          `json:"created_at"`
// 	Updated_at    time.Time          `json:"updated_at"`
// 	User_id       int                `json:"user_id"`
// }


type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	FirstName     string             `json:"first_name" validate:"required,min=2,max=100"`
	LastName      string             `json:"last_name" validate:"required,min=2,max=100"`
	Password      string             `json:"password" validate:"required,min=6"`
	Email         string             `json:"email" validate:"required,email"`
	Phone         string             `json:"phone" validate:"required"`
	Token         *string            `json:"token,omitempty"` // Token and Refresh_token as pointers
	Refresh_token *string            `json:"refresh_token,omitempty"`
	User_type     string             `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"` // Store User_id as a string to match primitive.ObjectID
}

