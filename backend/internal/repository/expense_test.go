package repository_test

import (
	"testing"
	"time"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpenseRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "expenses")
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewExpenseRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	expense := &model.Expense{
		Base:     model.Base{UserID: userID},
		Name:     "Grocery Shopping",
		Amount:   decimal.NewFromFloat(1500),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
		Type:     model.ExpenseTypeExpense,
	}

	err := repo.Create(expense)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, expense.ID)
}

func TestExpenseRepository_ListByYear(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "expenses")
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewExpenseRepository(db)
	userID := uuid.New()

	dates := []time.Time{
		time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
	}

	for i, d := range dates {
		expense := &model.Expense{
			Base:     model.Base{UserID: userID},
			Name:     "Expense",
			Amount:   decimal.NewFromInt(int64(i+1) * 1000),
			Currency: "MXN",
			Date:     d,
			Year:     d.Year(),
			Type:     model.ExpenseTypeExpense,
		}
		require.NoError(t, repo.Create(expense))
	}

	list2026, err := repo.ListByYear(userID, 2026)
	require.NoError(t, err)
	assert.Len(t, list2026, 2)

	list2025, err := repo.ListByYear(userID, 2025)
	require.NoError(t, err)
	assert.Len(t, list2025, 1)
}

func TestExpenseRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "expenses")
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewExpenseRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	expense := &model.Expense{
		Base:     model.Base{UserID: userID},
		Name:     "Old Name",
		Amount:   decimal.NewFromFloat(5000),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
		Type:     model.ExpenseTypeExpense,
	}
	require.NoError(t, repo.Create(expense))

	expense.Name = "New Name"
	err := repo.Update(expense)
	require.NoError(t, err)

	fetched, err := repo.GetByID(userID, expense.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Name", fetched.Name)
}

func TestExpenseRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "expenses")
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewExpenseRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	expense := &model.Expense{
		Base:     model.Base{UserID: userID},
		Name:     "To Delete",
		Amount:   decimal.NewFromFloat(1000),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
		Type:     model.ExpenseTypeExpense,
	}
	require.NoError(t, repo.Create(expense))

	err := repo.Delete(userID, expense.ID)
	require.NoError(t, err)

	list, err := repo.ListByYear(userID, 2026)
	require.NoError(t, err)
	assert.Len(t, list, 0)
}
