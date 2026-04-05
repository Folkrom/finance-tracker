package main

import (
	"log"
	"os"

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

	db, err := database.New(cfg, logger)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Repositories
	categoryRepo := repository.NewCategoryRepository(db)
	paymentMethodRepo := repository.NewPaymentMethodRepository(db)
	incomeRepo := repository.NewIncomeRepository(db)

	// Services
	categorySvc := service.NewCategoryService(categoryRepo)
	paymentMethodSvc := service.NewPaymentMethodService(paymentMethodRepo)
	incomeSvc := service.NewIncomeService(incomeRepo)

	// Handlers
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	paymentMethodHandler := handler.NewPaymentMethodHandler(paymentMethodSvc)
	incomeHandler := handler.NewIncomeHandler(incomeSvc)

	// Fiber app
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	router.Setup(app, cfg.SupabaseJWTSecret, categoryHandler, paymentMethodHandler, incomeHandler)

	logger.Info("server starting", zap.String("port", cfg.Port))
	log.Fatal(app.Listen(":" + cfg.Port))
}
