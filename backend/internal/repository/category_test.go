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
		Base:   model.Base{},
		UserID: &userID,
		Name:   "Salary",
		Domain: model.CategoryDomainIncome,
	}

	err := repo.Create(cat)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, cat.ID)
}

func TestCategoryRepository_ListByDomain_IncludesGlobal(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	// Create a global category (user_id = nil)
	global := &model.Category{
		Name:   "Global Salary",
		Domain: model.CategoryDomainIncome,
	}
	require.NoError(t, repo.Create(global))

	// Create a user category
	user := &model.Category{
		UserID: &userID,
		Name:   "My Side Gig",
		Domain: model.CategoryDomainIncome,
	}
	require.NoError(t, repo.Create(user))

	// User should see both global and their own
	cats, err := repo.ListByDomain(userID, model.CategoryDomainIncome)
	require.NoError(t, err)
	assert.Len(t, cats, 2)

	// Other user should see only global
	other, err := repo.ListByDomain(uuid.New(), model.CategoryDomainIncome)
	require.NoError(t, err)
	assert.Len(t, other, 1)
	assert.Equal(t, "Global Salary", other[0].Name)
}

func TestCategoryRepository_GetByID_GlobalCategory(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)

	global := &model.Category{
		Name:   "Global Expense",
		Domain: model.CategoryDomainExpense,
	}
	require.NoError(t, repo.Create(global))

	// Any user can read a global category
	cat, err := repo.GetByID(uuid.New(), global.ID)
	require.NoError(t, err)
	assert.Equal(t, "Global Expense", cat.Name)
}

func TestCategoryRepository_Delete_CannotDeleteGlobal(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)

	global := &model.Category{
		Name:   "Undeletable",
		Domain: model.CategoryDomainExpense,
	}
	require.NoError(t, repo.Create(global))

	// Attempting to delete a global category with any user_id should not delete it
	err := repo.Delete(uuid.New(), global.ID)
	require.NoError(t, err)

	// Global category should still exist
	cat, err := repo.GetByID(uuid.New(), global.ID)
	require.NoError(t, err)
	assert.Equal(t, "Undeletable", cat.Name)
}

func TestCategoryRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "categories")

	repo := repository.NewCategoryRepository(db)
	userID := uuid.New()

	cat := &model.Category{
		UserID: &userID,
		Name:   "Old Name",
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
