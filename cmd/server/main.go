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
	linkHandler := handler.Link{
		LinkService: link.NewService(),
	}

	application := web.New()

	application.Method("POST", "/link", linkHandler.CreateLink())
	application.Method("GET", "/link/{id}", linkHandler.GetLink())

	return application.Run()
}
