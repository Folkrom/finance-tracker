package service

import (
	"errors"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrGlobalCategoryReadOnly  = errors.New("cannot modify global categories")
	ErrSystemCategoryProtected = errors.New("cannot delete system categories")
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(userID uuid.UUID, name string, domain model.CategoryDomain, color *string) (*model.Category, error) {
	cat := &model.Category{
		UserID: &userID,
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
	if cat.IsGlobal {
		return nil, ErrGlobalCategoryReadOnly
	}
	cat.Name = name
	cat.Color = color
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(userID, id uuid.UUID) error {
	cat, err := s.repo.GetByID(userID, id)
	if err != nil {
		return err
	}
	if cat.IsSystem {
		return ErrSystemCategoryProtected
	}
	if cat.IsGlobal {
		return ErrGlobalCategoryReadOnly
	}
	return s.repo.Delete(userID, id)
}

// SeedDefaults is deprecated — global categories exist from migration.
func (s *CategoryService) SeedDefaults(_ uuid.UUID) error {
	return nil
}

func (s *CategoryService) CreateGlobal(name string, domain model.CategoryDomain, color *string, sortOrder int) (*model.Category, error) {
	cat := &model.Category{
		Name:      name,
		Domain:    domain,
		Color:     color,
		SortOrder: sortOrder,
	}
	if err := s.repo.CreateGlobal(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) UpdateGlobal(id uuid.UUID, name string, color *string, sortOrder *int) (*model.Category, error) {
	cat, err := s.repo.GetGlobalByID(id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		cat.Name = name
	}
	if color != nil {
		cat.Color = color
	}
	if sortOrder != nil {
		cat.SortOrder = *sortOrder
	}
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) DeleteGlobal(id uuid.UUID) error {
	cat, err := s.repo.GetGlobalByID(id)
	if err != nil {
		return err
	}
	if cat.IsSystem {
		return ErrSystemCategoryProtected
	}
	other, err := s.repo.GetOtherCategory(cat.Domain)
	if err != nil {
		return err
	}
	return s.repo.ReassignAndDelete(cat.ID, other.ID)
}
