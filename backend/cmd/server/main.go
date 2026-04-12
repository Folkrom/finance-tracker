package main

import (
	"log"
	"os"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/folkrom/finance-tracker/backend/internal/config"
	"github.com/folkrom/finance-tracker/backend/internal/database"
	"github.com/folkrom/finance-tracker/backend/internal/handler"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/internal/router"
	"github.com/folkrom/finance-tracker/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	if os.Getenv("ENVIRONMENT") == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	jwksURL := cfg.SupabaseURL + "/auth/v1/.well-known/jwks.json"
	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		log.Fatalf("failed to fetch JWKS from %s: %v", jwksURL, err)
	}

	db, err := database.New(cfg, logger)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Repositories
	categoryRepo := repository.NewCategoryRepository(db)
	paymentMethodRepo := repository.NewPaymentMethodRepository(db)
	incomeRepo := repository.NewIncomeRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	debtRepo := repository.NewDebtRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	cardRepo := repository.NewCardRepository(db)
	wishlistRepo := repository.NewWishlistItemRepository(db)

	// Services
	categorySvc := service.NewCategoryService(categoryRepo)
	paymentMethodSvc := service.NewPaymentMethodService(paymentMethodRepo)
	incomeSvc := service.NewIncomeService(incomeRepo)
	expenseSvc := service.NewExpenseService(expenseRepo)
	debtSvc := service.NewDebtService(debtRepo)
	budgetSvc := service.NewBudgetService(budgetRepo, expenseRepo, debtRepo)
	dashboardSvc := service.NewDashboardService(incomeRepo, expenseRepo, debtRepo)
	cardSvc := service.NewCardService(cardRepo, debtRepo, paymentMethodRepo)
	wishlistSvc := service.NewWishlistItemService(wishlistRepo)

	// Handlers
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	paymentMethodHandler := handler.NewPaymentMethodHandler(paymentMethodSvc)
	incomeHandler := handler.NewIncomeHandler(incomeSvc)
	expenseHandler := handler.NewExpenseHandler(expenseSvc)
	debtHandler := handler.NewDebtHandler(debtSvc)
	budgetHandler := handler.NewBudgetHandler(budgetSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)
	cardHandler := handler.NewCardHandler(cardSvc)
	wishlistHandler := handler.NewWishlistItemHandler(wishlistSvc)

	// Fiber app
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	router.Setup(app, jwks.Keyfunc, categoryHandler, paymentMethodHandler, incomeHandler, expenseHandler, debtHandler, budgetHandler, dashboardHandler, cardHandler, wishlistHandler)

	logger.Info("server starting", zap.String("port", cfg.Port))
	log.Fatal(app.Listen(":" + cfg.Port))
}
