package main

import (
	"go-apis/mgo"
	"go-apis/routers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	mgo.InitMongoDB()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is running here!")
	})

	routers.MountRoutes(e)
	e.Logger.Fatal(e.Start(":8800"))
}
