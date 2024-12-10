package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ObjectID  primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Mobile    string             `bson:"phone"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
}

type CreateUserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Mobile    string             `json:"phone"`
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Mobile    string `json:"phone" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
