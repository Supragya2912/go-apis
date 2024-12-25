package routers

import (
	"go-apis/api/auth"
	"go-apis/api/user"
	"go-apis/middleware"

	"github.com/labstack/echo/v4"
)

func MountRoutes(e *echo.Echo) {
	e.POST("/api/create", user.Create)
	e.POST("/api/login", auth.Login)

	protected := e.Group("/api")
	protected.Use(middleware.Protect)
	protected.POST("/update-password", auth.UpdatePassword)
}
