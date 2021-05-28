package main

import (
	"log"

	"github.com/emacampolo/link-tracker/cmd/server/handler"
	"github.com/emacampolo/link-tracker/internal/link"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	linkRepository := link.NewInMemoryRepository()
	linkService := link.NewService(linkRepository)
	linkHandler := handler.NewLink(linkService)

	engine := gin.Default()
	engine.POST("/link", linkHandler.Create)
	engine.GET("/link/{id}", linkHandler.Redirect)
	engine.GET("/link/{id}/metrics", linkHandler.Metrics)
	engine.POST("/link/{id}/inactivate", linkHandler.Inactivate)

	return engine.Run()
}
