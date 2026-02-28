package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// NotificationRepository 基于发件箱模式提供通知持久化操作。
type NotificationRepository struct{ db *gorm.DB }

// NewNotificationRepository 创建通知仓储实例。
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateOutbox 将一条状态为 "pending" 的通知持久化到发件箱表。
// 应在业务事务内调用，以保证消息的原子性。
func (r *NotificationRepository) CreateOutbox(ctx context.Context, n *domain.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

// ListPending 返回最多 limit 条尚未投递的待处理通知。
func (r *NotificationRepository) ListPending(ctx context.Context, limit int) ([]domain.Notification, error) {
	var list []domain.Notification
	err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// MarkSent 将通知状态标记为 "sent"（已发送）。
func (r *NotificationRepository) MarkSent(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("id = ?", id).
		Update("status", "sent").Error
}

// MarkFailed 将通知状态标记为 "failed"（发送失败）。
func (r *NotificationRepository) MarkFailed(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("id = ?", id).
		Update("status", "failed").Error
}
