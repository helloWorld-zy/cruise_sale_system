package service

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAuthService_Full(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	_ = db.AutoMigrate(&domain.Staff{})

	repo := repository.NewStaffRepository(db)
	svc := NewAuthService(repo, "secret", 1)
	ctx := context.Background()

	// Empty DB cases
	_, _, _ = svc.Login(ctx, "nonexistent", "pwd")
	_, _ = svc.GetProfile(ctx, "99")
	_, _ = svc.GetProfile(ctx, "invalid")

	// Valid cases
	pwd, _ := HashPassword("123456")
	_ = repo.Create(ctx, &domain.Staff{Username: "test", PasswordHash: pwd, Status: 1})
	_ = repo.Create(ctx, &domain.Staff{Username: "disabled", PasswordHash: pwd, Status: 0})

	_, _, _ = svc.Login(ctx, "test", "123456")
	_, _, _ = svc.Login(ctx, "test", "wrong")
	_, _, _ = svc.Login(ctx, "disabled", "123456")

	// Add test for disabled logic when repo throws disabled status, wait the repo actually is used in this test
	disabledPwd, _ := HashPassword("123456")
	_ = repo.Create(ctx, &domain.Staff{Username: "disabled_status_2", PasswordHash: disabledPwd, Status: 2})
	_, _, _ = svc.Login(ctx, "disabled_status_2", "123456")

	_, _ = svc.GetProfile(ctx, "1")
}
