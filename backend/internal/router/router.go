package router

import (
	"github.com/folkrom/finance-tracker/backend/internal/handler"
	"github.com/folkrom/finance-tracker/backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	jwtSecret string,
	categoryHandler *handler.CategoryHandler,
	paymentMethodHandler *handler.PaymentMethodHandler,
	incomeHandler *handler.IncomeHandler,
) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1", middleware.NewAuthMiddleware(jwtSecret))

	// Categories
	categories := api.Group("/categories")
	categories.Post("/", categoryHandler.Create)
	categories.Get("/", categoryHandler.List)
	categories.Post("/seed", categoryHandler.SeedDefaults)
	categories.Get("/:id", categoryHandler.GetByID)
	categories.Put("/:id", categoryHandler.Update)
	categories.Delete("/:id", categoryHandler.Delete)

	// Payment Methods
	paymentMethods := api.Group("/payment-methods")
	paymentMethods.Post("/", paymentMethodHandler.Create)
	paymentMethods.Get("/", paymentMethodHandler.List)
	paymentMethods.Get("/:id", paymentMethodHandler.GetByID)
	paymentMethods.Put("/:id", paymentMethodHandler.Update)
	paymentMethods.Delete("/:id", paymentMethodHandler.Delete)

	// Income (year-scoped)
	api.Post("/years/:year/incomes", incomeHandler.Create)
	api.Get("/years/:year/incomes", incomeHandler.ListByYear)
	api.Get("/incomes/:id", incomeHandler.GetByID)
	api.Put("/incomes/:id", incomeHandler.Update)
	api.Delete("/incomes/:id", incomeHandler.Delete)
}
