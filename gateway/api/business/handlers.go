package business

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"modular_chassis/gateway/internal/service/swagger"
	"modular_chassis/mediator/pkg/api"
	"reflect"
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

	reqModel, _ := api.GetMediatorAPI().GetServiceModels(domain, service)

	request := reflect.New(reqModel).Interface()
	err := json.Unmarshal(c.Body(), request)
	if err != nil {
		return err
	}

	response, err := api.GetMediatorAPI().Route(c.Context(), domain, service, request)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(response)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/json")
	return c.SendString(string(marshal))
}
