package main

import (
	"log"

	"link-tracker/cmd/server/handler"
	"link-tracker/internal/link"
	"link-tracker/internal/platform/web"
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

	application := web.New()

	application.Method("POST", "/link", linkHandler.Create())
	application.Method("GET", "/link/{id}", linkHandler.Redirect())

	return application.Run()
}
