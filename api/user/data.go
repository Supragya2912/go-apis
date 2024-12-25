package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type ExistingUserRequest struct {
	Email string `json:"email" validate:"required"`
}

type ExistingUserResponse struct {
	Exists bool `json:"exists"`
}
