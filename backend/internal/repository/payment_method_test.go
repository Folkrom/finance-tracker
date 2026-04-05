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

func TestPaymentMethodRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "payment_methods")

	repo := repository.NewPaymentMethodRepository(db)
	userID := uuid.New()

	pm := &model.PaymentMethod{
		Base: model.Base{UserID: userID},
		Name: "My Cash",
		Type: model.PaymentMethodCash,
	}

	err := repo.Create(pm)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, pm.ID)
}

func TestPaymentMethodRepository_ListByUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "payment_methods")

	repo := repository.NewPaymentMethodRepository(db)
	userID := uuid.New()

	pms := []model.PaymentMethod{
		{Base: model.Base{UserID: userID}, Name: "Cash", Type: model.PaymentMethodCash},
		{Base: model.Base{UserID: userID}, Name: "BBVA Debit", Type: model.PaymentMethodDebitCard},
	}
	for i := range pms {
		require.NoError(t, repo.Create(&pms[i]))
	}

	list, err := repo.ListByUser(userID)
	require.NoError(t, err)
	assert.Len(t, list, 2)

	other, err := repo.ListByUser(uuid.New())
	require.NoError(t, err)
	assert.Len(t, other, 0)
}

func TestPaymentMethodRepository_ListByType(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "payment_methods")

	repo := repository.NewPaymentMethodRepository(db)
	userID := uuid.New()

	pms := []model.PaymentMethod{
		{Base: model.Base{UserID: userID}, Name: "Cash", Type: model.PaymentMethodCash},
		{Base: model.Base{UserID: userID}, Name: "BBVA Debit", Type: model.PaymentMethodDebitCard},
		{Base: model.Base{UserID: userID}, Name: "Santander Debit", Type: model.PaymentMethodDebitCard},
	}
	for i := range pms {
		require.NoError(t, repo.Create(&pms[i]))
	}

	debit, err := repo.ListByType(userID, model.PaymentMethodDebitCard)
	require.NoError(t, err)
	assert.Len(t, debit, 2)
}

func TestPaymentMethodRepository_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "payment_methods")

	repo := repository.NewPaymentMethodRepository(db)
	userID := uuid.New()

	pm := &model.PaymentMethod{
		Base: model.Base{UserID: userID},
		Name: "Old Name",
		Type: model.PaymentMethodCash,
	}
	require.NoError(t, repo.Create(pm))

	pm.Name = "New Name"
	err := repo.Update(pm)
	require.NoError(t, err)

	fetched, err := repo.GetByID(userID, pm.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Name", fetched.Name)
}

func TestPaymentMethodRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.CleanTable(t, db, "payment_methods")

	repo := repository.NewPaymentMethodRepository(db)
	userID := uuid.New()

	pm := &model.PaymentMethod{
		Base: model.Base{UserID: userID},
		Name: "To Delete",
		Type: model.PaymentMethodCash,
	}
	require.NoError(t, repo.Create(pm))

	err := repo.Delete(userID, pm.ID)
	require.NoError(t, err)

	list, err := repo.ListByUser(userID)
	require.NoError(t, err)
	assert.Len(t, list, 0)
}
