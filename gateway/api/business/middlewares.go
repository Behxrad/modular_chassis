package business

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"modular_chassis/gateway/internal/adaptor"
)

var panicInterceptor = func(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			var err error
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
			response, httpStatus := adaptor.GetErrHTTPAdaptor().MakeErroneousResponse(err)
			c.Status(httpStatus)
			c.JSON(response)
		}
	}()
	return c.Next()
}

var errorInterceptor = func(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		response, httpStatus := adaptor.GetErrHTTPAdaptor().MakeErroneousResponse(err)
		c.Status(httpStatus)
		return c.JSON(response)
	}
	return nil
}
