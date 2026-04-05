package repository_test

import (
	"testing"

	"github.com/folkrom/finance-tracker/backend/internal/model"
	"github.com/folkrom/finance-tracker/backend/internal/repository"
	"github.com/folkrom/finance-tracker/backend/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	cat := &model.Category{
		Base:  model.Base{UserID: userID},
		Name:  "Salary",
		Domain: model.CategoryDomainIncome,
	}

	err := repo.Create(cat)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, cat.ID)
}

func TestCategoryRepository_ListByDomain(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	cats := []model.Category{
		{Base: model.Base{UserID: userID}, Name: "Salary", Domain: model.CategoryDomainIncome},
		{Base: model.Base{UserID: userID}, Name: "Bonus", Domain: model.CategoryDomainIncome},
		{Base: model.Base{UserID: userID}, Name: "Groceries", Domain: model.CategoryDomainExpense},
	}
	for i := range cats {
		require.NoError(t, repo.Create(&cats[i]))
	}

	income, err := repo.ListByDomain(userID, model.CategoryDomainIncome)
	require.NoError(t, err)
	assert.Len(t, income, 2)

	expense, err := repo.ListByDomain(userID, model.CategoryDomainExpense)
	require.NoError(t, err)
	assert.Len(t, expense, 1)

	// Other user should see nothing
	other, err := repo.ListByDomain(uuid.New(), model.CategoryDomainIncome)
	require.NoError(t, err)
	assert.Len(t, other, 0)
}

func TestCategoryRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	cat := &model.Category{
		Base:  model.Base{UserID: userID},
		Name:  "Old Name",
		Domain: model.CategoryDomainIncome,
	}
	require.NoError(t, repo.Create(cat))

	cat.Name = "New Name"
	err := repo.Update(cat)
	require.NoError(t, err)

	fetched, err := repo.GetByID(userID, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Name", fetched.Name)
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	cat := &model.Category{
		Base:  model.Base{UserID: userID},
		Name:  "To Delete",
		Domain: model.CategoryDomainExpense,
	}
	require.NoError(t, repo.Create(cat))

	err := repo.Delete(userID, cat.ID)
	require.NoError(t, err)

	list, err := repo.ListByDomain(userID, model.CategoryDomainExpense)
	require.NoError(t, err)
	assert.Len(t, list, 0)
}
