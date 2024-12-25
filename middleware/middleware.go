package middleware

import (
	"context"
	"errors"
	"fmt"
	"go-apis/api/user"
	"go-apis/helpers"
	"go-apis/mgo"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type jwtClaims struct {
	UserID string `json:"_id" bson:"_id"`
	jwt.StandardClaims
}

type User struct {
	ObjectID  primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Mobile    string             `bson:"phone"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
}

func Protect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, helpers.ErrResponse(
				http.StatusUnauthorized,
				"Authorization header is required",
				"AUTH_HEADER_MISSING",
			))
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := verifyToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, helpers.ErrResponse(
				http.StatusUnauthorized,
				"Invalid access token",
				"INVALID_ACCESS_TOKEN",
			))
		}

		filter := bson.M{"email": claims.Subject}

		var user user.CreateUserRequest
		err = mgo.Users.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusUnauthorized, helpers.ErrResponse(
					http.StatusUnauthorized,
					"User not found",
					"USER_NOT_FOUND",
				))
			}
			return c.JSON(http.StatusInternalServerError, helpers.ErrResponse(
				http.StatusInternalServerError,
				"Failed to fetch user",
				"FETCH_USER_ERROR",
			))
		}

		c.Set("user", user)

		return next(c)
	}
}

func verifyToken(tokenString string) (*jwtClaims, error) {
	fmt.Println("Received Token:", tokenString)

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
