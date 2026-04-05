package service

import (
	"errors"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CardService struct {
	repo              *repository.CardRepository
	debtRepo          *repository.DebtRepository
	paymentMethodRepo *repository.PaymentMethodRepository
}

func NewCardService(repo *repository.CardRepository, debtRepo *repository.DebtRepository, paymentMethodRepo *repository.PaymentMethodRepository) *CardService {
	return &CardService{repo: repo, debtRepo: debtRepo, paymentMethodRepo: paymentMethodRepo}
}

type CardSummary struct {
	Card           model.Card `json:"card"`
	AutoUsage      string     `json:"auto_usage"`
	ManualOverride *string    `json:"manual_override,omitempty"`
	TotalUsage     string     `json:"total_usage"`
	UsagePercent   float64    `json:"usage_percent"`
	RecommendedMax string     `json:"recommended_max"`
	HealthColor    string     `json:"health_color"`
}

func (s *CardService) Create(userID uuid.UUID, paymentMethodID uuid.UUID, bank string, cardLimit, recommendedMaxPct decimal.Decimal, manualOverride *decimal.Decimal, level *string) (*model.Card, error) {
	pm, err := s.paymentMethodRepo.GetByID(userID, paymentMethodID)
	if err != nil {
		return nil, errors.New("payment method not found")
	}
	if pm.Type != model.PaymentMethodCreditCard {
		return nil, errors.New("payment method must be of type credit_card")
	}

	card := &model.Card{
		Base:                model.Base{UserID: userID},
		PaymentMethodID:     paymentMethodID,
		Bank:                bank,
		CardLimit:           cardLimit,
		RecommendedMaxPct:   recommendedMaxPct,
		ManualUsageOverride: manualOverride,
		Level:               level,
	}
	if err := s.repo.Create(card); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, card.ID)
}

func (s *CardService) GetCardSummaries(userID uuid.UUID, month, year int) ([]CardSummary, error) {
	cards, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	debtSums, err := s.debtRepo.SumByPaymentMethodMonth(userID, month, year)
	if err != nil {
		return nil, err
	}

	debtMap := make(map[string]decimal.Decimal)
	for _, ds := range debtSums {
		debtMap[ds.PaymentMethodID] = ds.Total
	}

	summaries := make([]CardSummary, len(cards))
	for i, card := range cards {
		autoUsage := debtMap[card.PaymentMethodID.String()]
		totalUsage := autoUsage
		var manualStr *string
		if card.ManualUsageOverride != nil {
			totalUsage = totalUsage.Add(*card.ManualUsageOverride)
			s := card.ManualUsageOverride.String()
			manualStr = &s
		}

		var usagePct float64
		if !card.CardLimit.IsZero() {
			usagePct, _ = totalUsage.Mul(decimal.NewFromInt(100)).Div(card.CardLimit).Float64()
		}

		recommendedMax := card.CardLimit.Mul(card.RecommendedMaxPct).Div(decimal.NewFromInt(100))

		healthColor := "green"
		if usagePct > 70 {
			healthColor = "red"
		} else if usagePct > 30 {
			healthColor = "orange"
		} else if usagePct > 20 {
			healthColor = "yellow"
		}

		summaries[i] = CardSummary{
			Card:           card,
			AutoUsage:      autoUsage.String(),
			ManualOverride: manualStr,
			TotalUsage:     totalUsage.String(),
			UsagePercent:   usagePct,
			RecommendedMax: recommendedMax.String(),
			HealthColor:    healthColor,
		}
	}

	return summaries, nil
}

func (s *CardService) GetByID(userID, id uuid.UUID) (*model.Card, error) {
	return s.repo.GetByID(userID, id)
}

func (s *CardService) Update(userID, id uuid.UUID, bank string, cardLimit, recommendedMaxPct decimal.Decimal, manualOverride *decimal.Decimal, level *string) (*model.Card, error) {
	card, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	card.Bank = bank
	card.CardLimit = cardLimit
	card.RecommendedMaxPct = recommendedMaxPct
	card.ManualUsageOverride = manualOverride
	card.Level = level

	if err := s.repo.Update(card); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, card.ID)
}

func (s *CardService) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
