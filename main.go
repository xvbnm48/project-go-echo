package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func yallo(c echo.Context) error {
	return c.String(http.StatusOK, "hello")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("you cat is %s and type is %s", catName, catType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "invalid data type",
	})
}

func addCat(c echo.Context) error {
	cat := Cat{}
	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("failed reading the request body for add cats: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("failed unmarshaling the request body for add cats: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your cat: %v", cat)
	return c.JSON(http.StatusOK, "we got your cat")
}

func addDog(c echo.Context) error {
	dog := Dog{}
	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("failed reading the request body for add dogs: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &dog)
	if err != nil {
		log.Printf("failed unmarshaling the request body for add dogs: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your dog: %v", dog)
	return c.JSON(http.StatusOK, "we got your dog")
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

	e.GET("/", yallo)
	e.GET("/cats/:data", getCats)
	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.Start(":8000")
}
