package service

import (
	"errors"
	"slices"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	allowedCurrencies = []string{"MXN", "USD", "EUR", "GBP", "BRL", "COP", "ARS"}
	allowedLanguages  = []string{"en", "es"}

	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidLanguage = errors.New("invalid language")
)

type ProfileService struct {
	repo *repository.ProfileRepository
}

func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetOrCreate(userID uuid.UUID) (*model.Profile, error) {
	profile, err := s.repo.GetByUserID(userID)
	if err == nil {
		return profile, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	profile = &model.Profile{
		Base:     model.Base{UserID: userID},
		Currency: "MXN",
		Language: "en",
	}
	if err := s.repo.Create(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *ProfileService) Update(userID uuid.UUID, currency, language string) (*model.Profile, error) {
	profile, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if currency != "" {
		if !slices.Contains(allowedCurrencies, currency) {
			return nil, ErrInvalidCurrency
		}
		profile.Currency = currency
	}

	if language != "" {
		if !slices.Contains(allowedLanguages, language) {
			return nil, ErrInvalidLanguage
		}
		profile.Language = language
	}

	if err := s.repo.Update(profile); err != nil {
		return nil, err
	}
	return profile, nil
}
