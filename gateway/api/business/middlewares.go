package business

import "github.com/gofiber/fiber/v2"

var panicInterceptor = func(c *fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {
			c.Status(fiber.StatusInternalServerError) //TODO: fix this later
		}
	}()
	return c.Next()
}
