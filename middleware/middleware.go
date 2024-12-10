package middleware

import (
	"context"
	"errors"
	"go-apis/helpers"
	"go-apis/mgo"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type jwtClaims struct {
	UserID string `json:"_id"`
	jwt.StandardClaims
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
		user := mgo.Users.FindOne(context.Background(), claims.UserID)

		if user.Err() != nil {
			if user.Err() == mongo.ErrNoDocuments {
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

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
