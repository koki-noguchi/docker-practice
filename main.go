package main

import (
	"fmt"
	"github.com/koki-noguchi/docker-practice/server"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/ws", server.HandleConnections)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	go server.HandleMessages()

	fmt.Println("Listening on port 8080")
	log.Fatal(e.Start(":8080"))
}
