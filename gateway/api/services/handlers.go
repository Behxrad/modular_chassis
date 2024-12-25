package services

import (
	"github.com/gofiber/fiber/v2"
	"modular_chassis/gateway/internal/service/swagger"
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
