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

func TestIncomeRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewIncomeRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	income := &model.Income{
		Base:     model.Base{UserID: userID},
		Source:   "Employer",
		Amount:   decimal.NewFromFloat(10000),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
	}

	err := repo.Create(income)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, income.ID)
}

func TestIncomeRepository_ListByYear(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewIncomeRepository(db)
	userID := uuid.New()

	dates := []time.Time{
		time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
	}

	for i, d := range dates {
		income := &model.Income{
			Base:     model.Base{UserID: userID},
			Source:   "Source",
			Amount:   decimal.NewFromInt(int64(i + 1) * 1000),
			Currency: "MXN",
			Date:     d,
			Year:     d.Year(),
		}
		require.NoError(t, repo.Create(income))
	}

	list2026, err := repo.ListByYear(userID, 2026)
	require.NoError(t, err)
	assert.Len(t, list2026, 2)

	list2025, err := repo.ListByYear(userID, 2025)
	require.NoError(t, err)
	assert.Len(t, list2025, 1)
}

func TestIncomeRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewIncomeRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	income := &model.Income{
		Base:     model.Base{UserID: userID},
		Source:   "Old Source",
		Amount:   decimal.NewFromFloat(5000),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
	}
	require.NoError(t, repo.Create(income))

	income.Source = "New Source"
	err := repo.Update(income)
	require.NoError(t, err)

	fetched, err := repo.GetByID(userID, income.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Source", fetched.Source)
}

func TestIncomeRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "incomes")
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewIncomeRepository(db)
	userID := uuid.New()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	income := &model.Income{
		Base:     model.Base{UserID: userID},
		Source:   "To Delete",
		Amount:   decimal.NewFromFloat(1000),
		Currency: "MXN",
		Date:     date,
		Year:     date.Year(),
	}
	require.NoError(t, repo.Create(income))

	err := repo.Delete(userID, income.ID)
	require.NoError(t, err)

	list, err := repo.ListByYear(userID, 2026)
	require.NoError(t, err)
	assert.Len(t, list, 0)
}
