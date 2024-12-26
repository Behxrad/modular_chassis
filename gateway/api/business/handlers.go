package business

import (
	"github.com/gofiber/fiber/v2"
	"modular_chassis/gateway/internal/service/swagger"
	"modular_chassis/mediator/pkg/api"
)

var swaggerJSONDoc = func(c *fiber.Ctx) error {
	swaggerJSON, err := swagger.GenerateSwagger()
	if err != nil {
		return err
	}
	err = c.SendString(swaggerJSON)
	if err != nil {
		return err
	}
	return nil
}

var generalAPI = func(c *fiber.Ctx) error {
	domain := c.Params("domain")
	service := c.Params("service")

	res, err := api.GetMediatorAPI().Route(c.Context(), domain, service, string(c.Body()))
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/json")
	return c.SendString(res)
}
