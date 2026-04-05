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

func TestDebtRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "debts")
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewDebtRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	debt := &model.Debt{
		Base:     model.Base{UserID: userID},
		Name:     "Credit Card",
		Amount:   decimal.NewFromFloat(3000),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
	}

	err := repo.Create(debt)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, debt.ID)
}

func TestDebtRepository_ListByYear(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "debts")
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewDebtRepository(db)
	userID := uuid.New()

	dates := []time.Time{
		time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
	}

	for i, d := range dates {
		debt := &model.Debt{
			Base:     model.Base{UserID: userID},
			Name:     "Debt",
			Amount:   decimal.NewFromInt(int64(i+1) * 1000),
			Currency: "MXN",
			Date:     d,
			Year:     d.Year(),
		}
		require.NoError(t, repo.Create(debt))
	}

	list2026, err := repo.ListByYear(userID, 2026)
	require.NoError(t, err)
	assert.Len(t, list2026, 2)

	list2025, err := repo.ListByYear(userID, 2025)
	require.NoError(t, err)
	assert.Len(t, list2025, 1)
}

func TestDebtRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "debts")
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewDebtRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	debt := &model.Debt{
		Base:     model.Base{UserID: userID},
		Name:     "Old Name",
		Amount:   decimal.NewFromFloat(5000),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
	}
	require.NoError(t, repo.Create(debt))

	debt.Name = "New Name"
	err := repo.Update(debt)
	require.NoError(t, err)

	fetched, err := repo.GetByID(userID, debt.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Name", fetched.Name)
}

func TestDebtRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "debts")
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewDebtRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	debt := &model.Debt{
		Base:     model.Base{UserID: userID},
		Name:     "To Delete",
		Amount:   decimal.NewFromFloat(1000),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
	}
	require.NoError(t, repo.Create(debt))

	err := repo.Delete(userID, debt.ID)
	require.NoError(t, err)

	list, err := repo.ListByYear(userID, 2026)
	require.NoError(t, err)
	assert.Len(t, list, 0)
}
