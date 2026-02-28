package repository

import (
	"context"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newNotificationTestRepo 创建通知仓储测试实例。
func newNotificationTestRepo(t *testing.T) *NotificationRepository {
	t.Helper()
	// 使用独立内存库，避免测试数据串扰。
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.Notification{}))
	return NewNotificationRepository(db)
}

func TestNotificationRepository_CreateOutbox(t *testing.T) {
	repo := newNotificationTestRepo(t)
	n := &domain.Notification{UserID: 1, Channel: "sms", Template: "booking_created", Payload: "{}", Status: "pending"}
	require.NoError(t, repo.CreateOutbox(context.Background(), n))
	assert.Greater(t, n.ID, int64(0))
}

func TestNotificationRepository_ListPending(t *testing.T) {
	repo := newNotificationTestRepo(t)
	db := repo.db
	now := time.Now()

	// 按创建时间插入多条 pending，验证升序返回。
	require.NoError(t, db.Create(&domain.Notification{UserID: 1, Channel: "sms", Template: "a", Payload: "{}", Status: "pending", CreatedAt: now.Add(2 * time.Minute)}).Error)
	require.NoError(t, db.Create(&domain.Notification{UserID: 2, Channel: "sms", Template: "b", Payload: "{}", Status: "pending", CreatedAt: now}).Error)
	require.NoError(t, db.Create(&domain.Notification{UserID: 3, Channel: "sms", Template: "c", Payload: "{}", Status: "sent", CreatedAt: now.Add(3 * time.Minute)}).Error)

	list, err := repo.ListPending(context.Background(), 10)
	require.NoError(t, err)
	require.Len(t, list, 2)
	assert.LessOrEqual(t, list[0].CreatedAt.UnixNano(), list[1].CreatedAt.UnixNano())
	assert.Equal(t, "pending", list[0].Status)
	assert.Equal(t, "pending", list[1].Status)
}

func TestNotificationRepository_ListPending_Limit(t *testing.T) {
	repo := newNotificationTestRepo(t)
	db := repo.db

	// 插入 3 条 pending，limit=2 时只返回两条。
	for i := 0; i < 3; i++ {
		require.NoError(t, db.Create(&domain.Notification{UserID: int64(i + 1), Channel: "sms", Template: "t", Payload: "{}", Status: "pending", CreatedAt: time.Now().Add(time.Duration(i) * time.Second)}).Error)
	}

	list, err := repo.ListPending(context.Background(), 2)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestNotificationRepository_MarkSent(t *testing.T) {
	repo := newNotificationTestRepo(t)
	n := &domain.Notification{UserID: 1, Channel: "sms", Template: "t", Payload: "{}", Status: "pending"}
	require.NoError(t, repo.CreateOutbox(context.Background(), n))
	require.NoError(t, repo.MarkSent(context.Background(), n.ID))

	list, err := repo.ListPending(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestNotificationRepository_MarkFailed(t *testing.T) {
	repo := newNotificationTestRepo(t)
	n := &domain.Notification{UserID: 1, Channel: "sms", Template: "t", Payload: "{}", Status: "pending"}
	require.NoError(t, repo.CreateOutbox(context.Background(), n))
	require.NoError(t, repo.MarkFailed(context.Background(), n.ID))

	list, err := repo.ListPending(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, list, 0)
}
