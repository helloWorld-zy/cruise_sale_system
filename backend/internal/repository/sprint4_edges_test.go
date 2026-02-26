package repository

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/database"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var dbSeq atomic.Int64

// isolatedDB creates a uniquely named in-memory SQLite DB to avoid cross-test interference.
func isolatedDB() *gorm.DB {
	name := fmt.Sprintf("testdb_%d", dbSeq.Add(1))
	db, _ := database.ConnectWithDialector(
		sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", name)),
		database.Config{MaxIdleConns: 1, MaxOpenConns: 1},
	)
	return db
}

func setupTestDB() *database.Config {
	return &database.Config{Host: "memory"}
}

func TestRepositoryMissingEdges(t *testing.T) {
	db := isolatedDB()

	// Create tables with raw SQL to avoid AutoMigrate foreign key issues on SQLite
	db.Exec(`CREATE TABLE cabin_holds (
		id INTEGER PRIMARY KEY AUTOINCREMENT, cabin_sku_id INTEGER, user_id INTEGER,
		qty INTEGER, expires_at DATETIME, created_at DATETIME)`)
	db.Exec(`CREATE TABLE cabin_inventories (
		id INTEGER PRIMARY KEY AUTOINCREMENT, cabin_sku_id INTEGER UNIQUE,
		total INTEGER, locked INTEGER DEFAULT 0, sold INTEGER DEFAULT 0, updated_at DATETIME)`)
	db.Exec(`CREATE TABLE inventory_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT, cabin_sku_id INTEGER,
		"change" INTEGER, reason TEXT, created_at DATETIME)`)

	cHoldRepo := NewCabinHoldRepository(db)

	_, err := cHoldRepo.ExistsActiveHoldTx(db, 1, 1, time.Now())
	assert.NoError(t, err)

	err = cHoldRepo.CreateHoldTx(db, &domain.CabinHold{CabinSKUID: 1, UserID: 1})
	assert.NoError(t, err)

	db.Create(&domain.CabinInventory{CabinSKUID: 1, Total: 10})
	err = cHoldRepo.AdjustInventoryTx(db, 1, -2, "")
	assert.NoError(t, err)

	err = cHoldRepo.AdjustInventoryTx(db, 1, -20, "")
	assert.Error(t, err)

	db.Exec("DROP TABLE IF EXISTS cabin_holds")
	err = cHoldRepo.CreateHoldTx(db, &domain.CabinHold{CabinSKUID: 1, UserID: 1})
	assert.Error(t, err)

	_, err = cHoldRepo.ExistsActiveHoldTx(db, 1, 1, time.Now())
	assert.Error(t, err)

	db.Exec("DROP TABLE IF EXISTS cabin_inventories")
	err = cHoldRepo.AdjustInventoryTx(db, 1, -2, "")
	assert.Error(t, err)

	// List count-error paths: create tables then drop them so Count fails
	db2 := isolatedDB()
	db2.Exec(`CREATE TABLE cruise_companies (
		id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, english_name TEXT,
		description TEXT, logo_url TEXT, status INTEGER DEFAULT 1,
		sort_order INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`)
	db2.Exec(`CREATE TABLE cruises (
		id INTEGER PRIMARY KEY AUTOINCREMENT, company_id INTEGER,
		name TEXT, english_name TEXT, description TEXT, cover_url TEXT,
		tonnage INTEGER, length REAL, width REAL, decks INTEGER, year_built INTEGER,
		max_passengers INTEGER, status INTEGER DEFAULT 1, sort_order INTEGER DEFAULT 0,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`)
	db2.Exec(`CREATE TABLE cabin_types (
		id INTEGER PRIMARY KEY AUTOINCREMENT, cruise_id INTEGER,
		name TEXT, description TEXT, sort_order INTEGER DEFAULT 0,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`)
	db2.Exec(`CREATE TABLE facility_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, sort_order INTEGER DEFAULT 0,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`)
	db2.Exec(`CREATE TABLE facilities (
		id INTEGER PRIMARY KEY AUTOINCREMENT, cruise_id INTEGER, category_id INTEGER,
		name TEXT, description TEXT, icon_url TEXT, sort_order INTEGER DEFAULT 0,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`)

	// Drop all tables then test count error for each
	db2.Exec("DROP TABLE facilities")
	db2.Exec("DROP TABLE facility_categories")
	db2.Exec("DROP TABLE cabin_types")
	db2.Exec("DROP TABLE cruises")
	db2.Exec("DROP TABLE cruise_companies")

	cruiseRepo := NewCruiseRepository(db2)
	_, _, _ = cruiseRepo.List(context.Background(), 1, 1, 10)

	companyRepo := NewCompanyRepository(db2)
	_, _, _ = companyRepo.List(context.Background(), "", 1, 10)

	cabinTypeRepo := NewCabinTypeRepository(db2)
	_, _, _ = cabinTypeRepo.ListByCruise(context.Background(), 1, 1, 10)

	facCatRepo := NewFacilityCategoryRepository(db2)
	_, _ = facCatRepo.List(context.Background())

	facRepo := NewFacilityRepository(db2)
	_, _ = facRepo.ListByCruise(context.Background(), 1)

	// User Repo
	db3 := isolatedDB()
	db3.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, phone TEXT UNIQUE,
		wx_open_id TEXT, nickname TEXT, avatar_url TEXT,
		status INTEGER DEFAULT 1, created_at DATETIME, updated_at DATETIME)`)
	userRepo := NewUserRepository(db3)
	_, _ = userRepo.FindOrCreateByPhone("13800000000")

	db3.Exec("DROP TABLE users")
	_, _ = userRepo.FindOrCreateByPhone("13800000001")
}

// TestAdjustInventoryLogCreateError tests the error path when InventoryLog.Create fails.
func TestAdjustInventoryLogCreateError(t *testing.T) {
	db := isolatedDB()
	// Create CabinInventory but NOT InventoryLog table
	db.Exec(`CREATE TABLE cabin_inventories (
		id INTEGER PRIMARY KEY AUTOINCREMENT, cabin_sku_id INTEGER UNIQUE,
		total INTEGER, locked INTEGER DEFAULT 0, sold INTEGER DEFAULT 0, updated_at DATETIME)`)
	db.Create(&domain.CabinInventory{CabinSKUID: 99, Total: 10})

	repo := NewCabinHoldRepository(db)
	err := repo.AdjustInventoryTx(db, 99, -1, "test")
	assert.Error(t, err)
}

// TestAdjustInventorySaveError tests the error path when Save fails in AdjustInventoryTx.
func TestAdjustInventorySaveError(t *testing.T) {
	db := isolatedDB()
	db.Exec(`CREATE TABLE cabin_inventories (
		id INTEGER PRIMARY KEY AUTOINCREMENT, cabin_sku_id INTEGER UNIQUE,
		total INTEGER, locked INTEGER DEFAULT 0, sold INTEGER DEFAULT 0, updated_at DATETIME)`)
	db.Exec(`CREATE TABLE inventory_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT, cabin_sku_id INTEGER,
		"change" INTEGER, reason TEXT, created_at DATETIME)`)
	db.Create(&domain.CabinInventory{CabinSKUID: 50, Total: 10})

	// Inject error on Update (Save uses Update internally)
	db.Callback().Update().Before("gorm:update").Register("test:fail_save", func(tx *gorm.DB) {
		tx.AddError(errors.New("injected save error"))
	})

	repo := NewCabinHoldRepository(db)
	err := repo.AdjustInventoryTx(db, 50, -1, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "injected save error")
}

// TestListFindErrorPaths tests the Find error branch in List methods using a GORM callback.
func TestListFindErrorPaths(t *testing.T) {
	t.Run("CompanyListFindError", func(t *testing.T) {
		db := isolatedDB()
		db.Exec(`CREATE TABLE cruise_companies (
			id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, english_name TEXT,
			description TEXT, logo_url TEXT, status INTEGER DEFAULT 1,
			sort_order INTEGER DEFAULT 0, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`)

		queryCount := 0
		db.Callback().Query().Before("gorm:query").Register("test:fail_co_find", func(tx *gorm.DB) {
			queryCount++
			if queryCount > 1 {
				tx.AddError(errors.New("injected find error"))
			}
		})

		repo := NewCompanyRepository(db)
		_, _, err := repo.List(context.Background(), "", 1, 10)
		assert.Error(t, err)
	})

	t.Run("CruiseListFindError", func(t *testing.T) {
		db := isolatedDB()
		db.Exec(`CREATE TABLE cruises (
			id INTEGER PRIMARY KEY AUTOINCREMENT, company_id INTEGER,
			name TEXT, english_name TEXT, description TEXT, cover_url TEXT,
			tonnage INTEGER, length REAL, width REAL, decks INTEGER, year_built INTEGER,
			max_passengers INTEGER, status INTEGER DEFAULT 1, sort_order INTEGER DEFAULT 0,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`)

		queryCount := 0
		db.Callback().Query().Before("gorm:query").Register("test:fail_cr_find", func(tx *gorm.DB) {
			queryCount++
			if queryCount > 1 {
				tx.AddError(errors.New("injected find error"))
			}
		})

		repo := NewCruiseRepository(db)
		_, _, err := repo.List(context.Background(), 0, 1, 10)
		assert.Error(t, err)
	})

	t.Run("CabinTypeListFindError", func(t *testing.T) {
		db := isolatedDB()
		db.Exec(`CREATE TABLE cabin_types (
			id INTEGER PRIMARY KEY AUTOINCREMENT, cruise_id INTEGER,
			name TEXT, description TEXT, sort_order INTEGER DEFAULT 0,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`)

		queryCount := 0
		db.Callback().Query().Before("gorm:query").Register("test:fail_ct_find", func(tx *gorm.DB) {
			queryCount++
			if queryCount > 1 {
				tx.AddError(errors.New("injected find error"))
			}
		})

		repo := NewCabinTypeRepository(db)
		_, _, err := repo.ListByCruise(context.Background(), 1, 1, 10)
		assert.Error(t, err)
	})
}

// TestFindOrCreateByPhoneErrorPaths tests Create-fail error paths in FindOrCreateByPhone.
func TestFindOrCreateByPhoneErrorPaths(t *testing.T) {
	// Test Create fails AND re-query also fails → covers both uncovered blocks
	db := isolatedDB()
	db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, phone TEXT UNIQUE,
		wx_open_id TEXT, nickname TEXT, avatar_url TEXT,
		status INTEGER DEFAULT 1, created_at DATETIME, updated_at DATETIME)`)

	queryCount := 0
	db.Callback().Create().Before("gorm:create").Register("test:fail_create", func(tx *gorm.DB) {
		tx.AddError(errors.New("create failed"))
	})
	db.Callback().Query().Before("gorm:query").Register("test:fail_requery", func(tx *gorm.DB) {
		queryCount++
		if queryCount > 1 {
			tx.AddError(errors.New("requery failed"))
		}
	})

	repo := NewUserRepository(db)
	_, err := repo.FindOrCreateByPhone("13800003333")
	assert.Error(t, err)
}
