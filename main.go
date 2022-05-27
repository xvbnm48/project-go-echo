package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("failed reading the request body for add dogs: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	log.Printf("this is your dog: %v", dog)
	return c.JSON(http.StatusOK, dog)
}

func addHamster(c echo.Context) error {
	hamster := Hamster{}
	defer c.Request().Body.Close()

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("failed reading the request body for add hamsters: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your hamster: %v", hamster)
	return c.JSON(http.StatusOK, hamster)
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "hello admin")
}
func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "you are on the not yet secret cookie page")
}

func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if username == "sakura" && password == "miyawaki" {
		// return c.JSON(http.StatusOK, map[string]string{
		// 	"message": "you are logged in",
		// })
		cookie := &http.Cookie{}

		// this is same
		// cookie := new(http.Cookie)
		cookie.Name = "sessionID"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)
		c.SetCookie(cookie)
		return c.JSON(http.StatusOK, map[string]string{
			"message": "you are logged in",
		})
	}
	return c.JSON(http.StatusUnauthorized, map[string]string{
		"message": "you are not logged in",
	})
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/4.0")
		c.Response().Header().Set("notReallyHeader", "thisHavenNoMeaning")

		return next(c)
	}
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")
		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.String(http.StatusUnauthorized, "you dont have any cookie")
			}

			log.Println(err)
			return err
		}

		if cookie.Value == "some_string" {
			return next(c)
		}

		return c.String(http.StatusUnauthorized, "you dont have the right cookie, cookie")
	}
}

func main() {
	port := os.Getenv("MY_APP_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Welcome to the server")

	e := echo.New()
	e.Use(ServerHeader)

	adminGroup := e.Group("/admin")
	cookieGroup := e.Group("/cookie")

	// this is logs the server interaction
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	adminGroup.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		// check in the DB
		if username == "sakura" && password == "miyawaki" {
			return true, nil
		}

		return false, nil
	}))

	cookieGroup.GET("/main", mainCookie)
	adminGroup.GET("/main", mainAdmin)

	e.GET("/login", login)
	e.GET("/", yallo)
	e.GET("/cats/:data", getCats)
	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/hamsters", addHamster)
	e.Logger.Print("listening on port " + port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%s", port)))
}
