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
	expenseHandler *handler.ExpenseHandler,
	debtHandler *handler.DebtHandler,
	budgetHandler *handler.BudgetHandler,
	dashboardHandler *handler.DashboardHandler,
	cardHandler *handler.CardHandler,
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

	// Expenses (year-scoped)
	api.Post("/years/:year/expenses", expenseHandler.Create)
	api.Get("/years/:year/expenses", expenseHandler.ListByYear)
	api.Get("/expenses/:id", expenseHandler.GetByID)
	api.Put("/expenses/:id", expenseHandler.Update)
	api.Delete("/expenses/:id", expenseHandler.Delete)

	// Debts (year-scoped)
	api.Post("/years/:year/debts", debtHandler.Create)
	api.Get("/years/:year/debts", debtHandler.ListByYear)
	api.Get("/debts/:id", debtHandler.GetByID)
	api.Put("/debts/:id", debtHandler.Update)
	api.Delete("/debts/:id", debtHandler.Delete)

	// Budgets
	budgets := api.Group("/budgets")
	budgets.Post("/", budgetHandler.Create)
	budgets.Get("/", budgetHandler.GetSummary)
	budgets.Get("/recurring", budgetHandler.ListRecurring)
	budgets.Put("/:id", budgetHandler.Update)
	budgets.Delete("/:id", budgetHandler.Delete)

	// Dashboard
	api.Get("/years/:year/dashboard", dashboardHandler.GetDashboard)

	// Cards
	cards := api.Group("/cards")
	cards.Post("/", cardHandler.Create)
	cards.Get("/", cardHandler.GetSummaries)
	cards.Get("/:id", cardHandler.GetByID)
	cards.Put("/:id", cardHandler.Update)
	cards.Delete("/:id", cardHandler.Delete)
}
