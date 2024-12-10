package routers

import (
	"go-apis/api/user"

	"github.com/labstack/echo/v4"
)

func MountRoutes(e *echo.Echo) {
	e.POST("/api/create", user.Create)
}
