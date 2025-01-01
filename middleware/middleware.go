package middleware

import (
	"context"

	"go-apis/api/user"
	"go-apis/helpers"
	"go-apis/mgo"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

		claims, err := helpers.VerifyToken(tokenString)
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
