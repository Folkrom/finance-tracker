package middleware

import (
	"sync"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// NewProfileMiddleware ensures a profile exists for every authenticated user.
func NewProfileMiddleware(profileRepo *repository.ProfileRepository) fiber.Handler {
	var seen sync.Map

	return func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		userIDStr := userID.String()

		if _, ok := seen.Load(userIDStr); ok {
			return c.Next()
		}

		_, err := profileRepo.GetByUserID(userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				profile := &model.Profile{
					Base:     model.Base{UserID: userID},
					Currency: "MXN",
					Language: "en",
				}
				if createErr := profileRepo.Create(profile); createErr != nil {
					if _, retryErr := profileRepo.GetByUserID(userID); retryErr != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": "failed to create profile",
						})
					}
				}
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to check profile",
				})
			}
		}

		seen.Store(userIDStr, true)
		return c.Next()
	}
}
