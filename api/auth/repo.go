package auth

import (
	"context"
	"errors"
	"go-apis/helpers"
	"go-apis/mgo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
	// Add more fields as needed
}

type Mongo interface {
	LoginUser(data *LoginRequest) (*LoginResponse, error)
}

type DefaultMongo struct{}

var dmgo = &DefaultMongo{}

func (d *DefaultMongo) LoginUser(data *LoginRequest) (*LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := mgo.Users.FindOne(ctx, bson.M{"email": data.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return nil, errors.New("invalid email or password")
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := helpers.GenerateAccessToken(user.Email)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
	}, nil
}
