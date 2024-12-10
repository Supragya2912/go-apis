package user

import (
	"fmt"
	"go-apis/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
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
