package middleware

import (
	"context"
	"go-apis/api/user"
	"go-apis/helpers"
	"go-apis/mgo"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckPermission(requiredRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userID, ok := c.Get("userID").(string)
			if !ok || userID == "" {
				return c.JSON(http.StatusUnauthorized, helpers.ErrResponse(
					http.StatusUnauthorized,
					"Unauthorized",
					"UNAUTHORIZED",
				))
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var user user.CreateUserRequest
			err := mgo.Users.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)

			if err != nil {
				if err == mongo.ErrNoDocuments {
					return c.JSON(http.StatusNotFound, helpers.ErrResponse(
						http.StatusNotFound,
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

			hasPermission := false
			for _, role := range requiredRoles {
				if user.Role == role {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				return c.JSON(http.StatusForbidden, helpers.ErrResponse(
					http.StatusInternalServerError,
					"You are unauthorized to perform this action",
					"UNAUTHORIZED_ACTION",
				))
			}

			return next(c)
		}
	}
}
