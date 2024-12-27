package business

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"sync"
)

var (
	httpServerOnce sync.Once
	httpServerIns  *httpServer
)

type httpServer struct {
}

func GetHTTPServer() *httpServer {
	httpServerOnce.Do(func() {
		if httpServerIns == nil {
			httpServerIns = &httpServer{}
		}
	})
	return httpServerIns
}

func (h httpServer) Run() error {
	app := fiber.New()

	app.Use(panicInterceptor)
	app.Use(errorInterceptor)

	app.Get("/swagger/doc.json", swaggerJSONDoc)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Post("/api/:domain/:service", generalAPI)

	err := app.Listen(":1323")
	if err != nil {
		return err
	}
	return nil
}
