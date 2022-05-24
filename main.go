package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func yallo(c echo.Context) error {
	return c.String(http.StatusOK, "hello")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	return c.String(http.StatusOK, fmt.Sprintf("you cat is %s and type is %s", catName, catType))
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

	e.GET("/", yallo)
	e.GET("/cats", getCats)

	e.Start(":8000")
}
