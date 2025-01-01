package user

import (
	"context"
	"fmt"

	"go-apis/helpers"
	"go-apis/mgo"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(c echo.Context) error {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			helpers.ErrResponse(
				http.StatusInternalServerError,
				"Failed to bind request",
				"BIND_ERROR",
			),
		)
	}

	// fmt.Println(req.Email, req.FirstName, req.LastName, req.Mobile, req.Password)

	if req.Email == "" || req.FirstName == "" || req.LastName == "" || req.Mobile == "" || req.Password == "" {
		return c.JSON(
			http.StatusBadRequest,
			helpers.ErrResponse(
				http.StatusBadRequest,
				"Missing required fields",
				"MISSING_FIELDS",
			),
		)
	}

	existingUserReq := &ExistingUserRequest{
		Email: req.Email,
	}
	exists, err := dmgo.ExistingUser(existingUserReq)
	fmt.Println(err, exists)

	if exists {
		return c.JSON(
			http.StatusUnauthorized,
			helpers.ErrResponse(
				http.StatusUnauthorized,
				"Email Exists",
				"USER_EXISTS",
			),
		)
	}

	req.Role = "user"
	user, err := dmgo.CreateUserRequest(req)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			helpers.ErrResponse(
				http.StatusInternalServerError,
				"Failed to create user",
				"CREATE_USER_ERROR",
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		helpers.SuccessResponse(user),
	)
}

func GetUser(c echo.Context) error {

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

	var user CreateUserRequest
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

	userResponse := GetUserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Mobile,
		Role:      user.Role,
	}
	return c.JSON(
		http.StatusOK,
		helpers.SuccessResponse(userResponse),
	)

}
