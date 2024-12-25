package services

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"sync"
)

var once sync.Once
var httpServerIns *httpServer

type httpServer struct {
}

func GetHTTPServer() *httpServer {
	once.Do(func() {
		if httpServerIns == nil {
			httpServerIns = &httpServer{}
		}
	})
	return httpServerIns
}

func (h httpServer) Run() error {
	app := fiber.New()

	app.Get("/swagger/doc.json", swaggerJSONDoc)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	err := app.Listen(":1323")
	if err != nil {
		return err
	}
	return nil
}