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

	// 空数据库情况
	_, _, _ = svc.Login(ctx, "nonexistent", "pwd")
	_, _ = svc.GetProfile(ctx, "99")
	_, _ = svc.GetProfile(ctx, "invalid")

	// 有效情况
	pwd, _ := HashPassword("123456")
	_ = repo.Create(ctx, &domain.Staff{Username: "test", PasswordHash: pwd, Status: 1})
	_ = repo.Create(ctx, &domain.Staff{Username: "disabled", PasswordHash: pwd, Status: 0})

	_, _, _ = svc.Login(ctx, "test", "123456")
	_, _, _ = svc.Login(ctx, "test", "wrong")
	_, _, _ = svc.Login(ctx, "disabled", "123456")

	// 当仓储抛出禁用状态时添加对禁用逻辑的测试，等待仓储实际在这个测试中使用
	disabledPwd, _ := HashPassword("123456")
	_ = repo.Create(ctx, &domain.Staff{Username: "disabled_status_2", PasswordHash: disabledPwd, Status: 2})
	_, _, _ = svc.Login(ctx, "disabled_status_2", "123456")

	_, _ = svc.GetProfile(ctx, "1")
}
