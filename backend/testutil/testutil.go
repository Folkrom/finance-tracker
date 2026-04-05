package testutil

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://finance:finance_dev@localhost:5432/finance_tracker?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return db
}

func CleanTable(t *testing.T, db *gorm.DB, table string) {
	t.Helper()
	db.Exec("DELETE FROM " + table)
}
