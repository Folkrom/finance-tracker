package handler

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Error string `json:"error"`
}

type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
}

func respondError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(ErrorResponse{Error: msg})
}

func respondList[T any](c *fiber.Ctx, data []T) error {
	if data == nil {
		data = []T{}
	}
	return c.JSON(ListResponse[T]{Data: data, Total: len(data)})
}

func respondJSON(c *fiber.Ctx, data interface{}) error {
	return c.JSON(data)
}

func respondCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

func respondNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
