package business

import (
	"encoding/json"
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

	request, err := api.GetMediatorAPI().GetServiceRequestModel(domain, service)
	if err != nil {
		return err
	}
	baseReq := api.GetMediatorAPI().GetBaseReqFromModel(request)
	if baseReq != nil {
		baseReq.Mobile = "09##***####" //TODO: make authentication here to fill headers
	}

	err = json.Unmarshal(c.Body(), request)
	if err != nil {
		return err
	}

	response, err := api.GetMediatorAPI().Route(c.Context(), domain, service, request)
	if err != nil {
		return err
	}

	return c.JSON(response)
}
