package service

import (
	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(userID uuid.UUID, name string, domain model.CategoryDomain, color *string) (*model.Category, error) {
	cat := &model.Category{
		Base:   model.Base{UserID: userID},
		Name:   name,
		Domain: domain,
		Color:  color,
	}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) ListByDomain(userID uuid.UUID, domain model.CategoryDomain) ([]model.Category, error) {
	return s.repo.ListByDomain(userID, domain)
}

func (s *CategoryService) GetByID(userID, id uuid.UUID) (*model.Category, error) {
	return s.repo.GetByID(userID, id)
}

func (s *CategoryService) Update(userID, id uuid.UUID, name string, color *string) (*model.Category, error) {
	cat, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	cat.Name = name
	cat.Color = color
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}

func (s *CategoryService) SeedDefaults(userID uuid.UUID) error {
	incomeCategories := []string{
		"Salary", "Bonus", "Freelance", "Dividends", "Interest", "Side Hustle",
	}
	expenseCategories := []string{
		"Home Expenses", "Eating Out", "Self Care", "Coffee/Drink", "Entertainment",
		"Transportation", "Groceries", "Utilities", "Clothes", "Other",
		"Card Payments", "Savings/Investment", "Alcohol", "Drugs", "Taxes",
		"Knowledge", "Tech",
	}

	for i, name := range incomeCategories {
		cat := &model.Category{
			Base:      model.Base{UserID: userID},
			Name:      name,
			Domain:    model.CategoryDomainIncome,
			SortOrder: i,
		}
		// Ignore duplicate errors (idempotent seeding)
		_ = s.repo.Create(cat)
	}

	for i, name := range expenseCategories {
		cat := &model.Category{
			Base:      model.Base{UserID: userID},
			Name:      name,
			Domain:    model.CategoryDomainExpense,
			SortOrder: i,
		}
		_ = s.repo.Create(cat)
	}

	wishlistCategories := []string{
		"Electronics", "Clothing", "Home & Kitchen", "Books & Media",
		"Sports & Outdoors", "Beauty & Personal Care", "Toys & Games", "Other",
	}
	for i, name := range wishlistCategories {
		cat := &model.Category{
			Base:      model.Base{UserID: userID},
			Name:      name,
			Domain:    model.CategoryDomainWishlist,
			SortOrder: i,
		}
		_ = s.repo.Create(cat)
	}

	return nil
}
