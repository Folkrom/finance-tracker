package repository_test

import (
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudgetRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "budgets")
	testutil.CleanTable(t, db, "categories")

	catRepo := repository.NewCategoryRepository(db)
	cat := &model.Category{
		Base: model.Base{UserID: uuid.New()},
		Name: "Food",
		Domain: model.CategoryDomainExpense,
	}
	require.NoError(t, catRepo.Create(cat))

	repo := repository.NewBudgetRepository(db)
	userID := uuid.New()

	budget := &model.Budget{
		Base:         model.Base{UserID: userID},
		CategoryID:   cat.ID,
		MonthlyLimit: decimal.NewFromFloat(5000),
		Month:        1,
		Year:         2026,
		IsRecurring:  false,
	}

	err := repo.Create(budget)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, budget.ID)
}

func TestBudgetRepository_ListByMonthYear(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "budgets")
	testutil.CleanTable(t, db, "categories")

	catRepo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	cat1 := &model.Category{Base: model.Base{UserID: userID}, Name: "Food", Domain: model.CategoryDomainExpense}
	cat2 := &model.Category{Base: model.Base{UserID: userID}, Name: "Transport", Domain: model.CategoryDomainExpense}
	require.NoError(t, catRepo.Create(cat1))
	require.NoError(t, catRepo.Create(cat2))

	repo := repository.NewBudgetRepository(db)

	// Specific override for Jan 2026
	override := &model.Budget{
		Base:         model.Base{UserID: userID},
		CategoryID:   cat1.ID,
		MonthlyLimit: decimal.NewFromFloat(3000),
		Month:        1,
		Year:         2026,
		IsRecurring:  false,
	}
	require.NoError(t, repo.Create(override))

	// Recurring template
	recurring := &model.Budget{
		Base:         model.Base{UserID: userID},
		CategoryID:   cat2.ID,
		MonthlyLimit: decimal.NewFromFloat(1500),
		Month:        0,
		Year:         0,
		IsRecurring:  true,
	}
	// For recurring, bypass unique constraint by inserting directly
	require.NoError(t, db.Exec(
		"INSERT INTO budgets (id, user_id, category_id, monthly_limit, month, year, is_recurring) VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, ?)",
		userID, cat2.ID, decimal.NewFromFloat(1500), 1, 2026, true,
	).Error)
	_ = recurring

	list, err := repo.ListByMonthYear(userID, 1, 2026)
	require.NoError(t, err)
	assert.Len(t, list, 2)

	// Different month should return only recurring
	list2, err := repo.ListByMonthYear(userID, 2, 2026)
	require.NoError(t, err)
	assert.Len(t, list2, 1)
}

func TestBudgetRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "budgets")
	testutil.CleanTable(t, db, "categories")

	catRepo := repository.NewCategoryRepository(db)
	cat := &model.Category{
		Base: model.Base{UserID: uuid.New()},
		Name: "Food",
		Domain: model.CategoryDomainExpense,
	}
	require.NoError(t, catRepo.Create(cat))

	repo := repository.NewBudgetRepository(db)
	userID := uuid.New()

	budget := &model.Budget{
		Base:         model.Base{UserID: userID},
		CategoryID:   cat.ID,
		MonthlyLimit: decimal.NewFromFloat(2000),
		Month:        3,
		Year:         2026,
		IsRecurring:  false,
	}
	require.NoError(t, repo.Create(budget))

	budget.MonthlyLimit = decimal.NewFromFloat(4000)
	err := repo.Update(budget)
	require.NoError(t, err)

	fetched, err := repo.GetByID(userID, budget.ID)
	require.NoError(t, err)
	assert.True(t, fetched.MonthlyLimit.Equal(decimal.NewFromFloat(4000)))
}

func TestBudgetRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "budgets")
	testutil.CleanTable(t, db, "categories")

	catRepo := repository.NewCategoryRepository(db)
	cat := &model.Category{
		Base: model.Base{UserID: uuid.New()},
		Name: "Food",
		Domain: model.CategoryDomainExpense,
	}
	require.NoError(t, catRepo.Create(cat))

	repo := repository.NewBudgetRepository(db)
	userID := uuid.New()

	budget := &model.Budget{
		Base:         model.Base{UserID: userID},
		CategoryID:   cat.ID,
		MonthlyLimit: decimal.NewFromFloat(1000),
		Month:        4,
		Year:         2026,
		IsRecurring:  false,
	}
	require.NoError(t, repo.Create(budget))

	err := repo.Delete(userID, budget.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(userID, budget.ID)
	assert.Error(t, err)
}
