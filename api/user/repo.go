package user

import (
	"context"
	"go-apis/mgo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/crypto/bcrypt"
)

type Mongo interface {
	CreateUserRequest(data *CreateUserRequest) (*mongo.InsertOneResult, error)
}

type DefaultMongo struct{}

var dmgo DefaultMongo

func (m DefaultMongo) CreateUserRequest(data *CreateUserRequest) (*CreateUserResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	data.Password = string(hashedPassword)

	insertResult, err := mgo.Users.InsertOne(context.Background(), data)
	if err != nil {
		return nil, err
	}

	var createdUser CreateUserResponse
	filter := bson.M{"_id": insertResult.InsertedID}
	err = mgo.Users.FindOne(context.Background(), filter).Decode(&createdUser)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}
