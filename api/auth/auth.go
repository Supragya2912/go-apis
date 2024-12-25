package auth

import (
	"fmt"
	"go-apis/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	req := new(LoginRequest)
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

	if req.Email == "" || req.Password == "" {
		return c.JSON(
			http.StatusBadRequest,
			helpers.ErrResponse(
				http.StatusBadRequest,
				"Missing required fields",
				"MISSING_FIELDS",
			),
		)
	}

	fmt.Print(req.Email, req.Password)

	token, err := dmgo.LoginUser(req)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			helpers.ErrResponse(
				http.StatusInternalServerError,
				err.Error(),
				"LOGIN_ERROR",
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		helpers.SuccessResponse(token),
	)
}

func UpdatePassword(c echo.Context) error {
	fmt.Println("Update Password Request Received")
	req :=
		new(UpdatePasswordRequest)
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

	if req.Email == "" || req.NewPassword == "" || req.OldPassword == "" {
		return c.JSON(
			http.StatusBadRequest,
			helpers.ErrResponse(
				http.StatusBadRequest,
				"Missing required fields",
				"MISSING_FIELDS",
			),
		)
	}

	err := dmgo.UpdatePassword(req)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			helpers.ErrResponse(
				http.StatusInternalServerError,
				err.Error(),
				"UPDATE_PASSWORD_ERROR",
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		helpers.SuccessResponse("Password updated successfully"),
	)

}
