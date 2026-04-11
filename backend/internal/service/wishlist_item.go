package service

import (
	"fmt"
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type WishlistItemService struct {
	repo *repository.WishlistItemRepository
}

func NewWishlistItemService(repo *repository.WishlistItemRepository) *WishlistItemService {
	return &WishlistItemService{repo: repo}
}

var validPriorities = map[model.WishlistPriority]bool{
	model.WishlistPriorityLow:    true,
	model.WishlistPriorityMedium: true,
	model.WishlistPriorityHigh:   true,
}

var validStatuses = map[model.WishlistStatus]bool{
	model.WishlistStatusInterested:     true,
	model.WishlistStatusSavingFor:      true,
	model.WishlistStatusWaitingForSale: true,
	model.WishlistStatusOrdered:        true,
	model.WishlistStatusPurchased:      true,
	model.WishlistStatusReceived:       true,
	model.WishlistStatusCancelled:      true,
}

func (s *WishlistItemService) Create(
	userID uuid.UUID,
	name string,
	imageURL *string,
	price *decimal.Decimal,
	currency string,
	links []string,
	categoryID *uuid.UUID,
	priority model.WishlistPriority,
	status model.WishlistStatus,
	targetDate *time.Time,
	monthlyContribution *decimal.Decimal,
	contributionCurrency string,
) (*model.WishlistItem, error) {
	if len(links) > 5 {
		return nil, fmt.Errorf("maximum 5 links allowed")
	}
	if !validPriorities[priority] {
		priority = model.WishlistPriorityMedium
	}
	if !validStatuses[status] {
		status = model.WishlistStatusInterested
	}
	if currency == "" {
		currency = "MXN"
	}
	if contributionCurrency == "" {
		contributionCurrency = "MXN"
	}

	item := &model.WishlistItem{
		Base:                 model.Base{UserID: userID},
		Name:                 name,
		ImageURL:             imageURL,
		Price:                price,
		Currency:             currency,
		Links:                pq.StringArray(links),
		CategoryID:           categoryID,
		Priority:             priority,
		Status:               status,
		TargetDate:           targetDate,
		MonthlyContribution:  monthlyContribution,
		ContributionCurrency: contributionCurrency,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, item.ID)
}

func (s *WishlistItemService) ListByUser(userID uuid.UUID) ([]model.WishlistItem, error) {
	return s.repo.ListByUser(userID)
}

func (s *WishlistItemService) ListByStatus(userID uuid.UUID, statuses []model.WishlistStatus) ([]model.WishlistItem, error) {
	return s.repo.ListByStatus(userID, statuses)
}

func (s *WishlistItemService) GetByID(userID, id uuid.UUID) (*model.WishlistItem, error) {
	return s.repo.GetByID(userID, id)
}

func (s *WishlistItemService) Update(
	userID, id uuid.UUID,
	name string,
	imageURL *string,
	price *decimal.Decimal,
	currency string,
	links []string,
	categoryID *uuid.UUID,
	priority model.WishlistPriority,
	status model.WishlistStatus,
	targetDate *time.Time,
	monthlyContribution *decimal.Decimal,
	contributionCurrency string,
) (*model.WishlistItem, error) {
	if len(links) > 5 {
		return nil, fmt.Errorf("maximum 5 links allowed")
	}
	if !validPriorities[priority] {
		priority = model.WishlistPriorityMedium
	}
	if !validStatuses[status] {
		status = model.WishlistStatusInterested
	}
	if currency == "" {
		currency = "MXN"
	}
	if contributionCurrency == "" {
		contributionCurrency = "MXN"
	}

	item, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}

	item.Name = name
	item.ImageURL = imageURL
	item.Price = price
	item.Currency = currency
	item.Links = pq.StringArray(links)
	item.CategoryID = categoryID
	item.Priority = priority
	item.Status = status
	item.TargetDate = targetDate
	item.MonthlyContribution = monthlyContribution
	item.ContributionCurrency = contributionCurrency

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return s.repo.GetByID(userID, item.ID)
}

func (s *WishlistItemService) UpdateStatus(userID, id uuid.UUID, status model.WishlistStatus) error {
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	return s.repo.UpdateStatus(userID, id, status)
}

func (s *WishlistItemService) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
