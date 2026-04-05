package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func parseMoneyFields(amountStr string, categoryIDStr *string, paymentMethodIDStr *string, dateStr string) (
	decimal.Decimal, *uuid.UUID, *uuid.UUID, time.Time, error,
) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return decimal.Zero, nil, nil, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid amount")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return decimal.Zero, nil, nil, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid date, expected YYYY-MM-DD")
	}

	var categoryID *uuid.UUID
	if categoryIDStr != nil && *categoryIDStr != "" {
		parsed, err := uuid.Parse(*categoryIDStr)
		if err != nil {
			return decimal.Zero, nil, nil, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid category_id")
		}
		categoryID = &parsed
	}

	var paymentMethodID *uuid.UUID
	if paymentMethodIDStr != nil && *paymentMethodIDStr != "" {
		parsed, err := uuid.Parse(*paymentMethodIDStr)
		if err != nil {
			return decimal.Zero, nil, nil, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid payment_method_id")
		}
		paymentMethodID = &parsed
	}

	return amount, categoryID, paymentMethodID, date, nil
}
