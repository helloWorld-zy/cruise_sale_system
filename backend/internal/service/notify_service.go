package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// 通知状态常量。
const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

// 通知渠道常量。
const (
	ChannelSMS    = "sms"
	ChannelWechat = "wechat"
	ChannelInbox  = "inbox"
)

// NotifyService 实现发件箱模式以实现可靠的通知交付。
//
// 与同步调用外部提供商（在发生暂时性故障时会丢失消息）不同，
// Enqueue 与调用者的业务事务原子地将待处理通知写入发件箱表。
// 随后，一个独立的 OutboxDispatcher goroutine 读取待处理记录，
// 调用真实的提供商，并将每条通知标记为已发送或失败，失败时重试。
type NotifyService struct {
	repo domain.NotificationRepository
}

// NewNotifyService 创建一个由给定仓储支持的 NotifyService。
func NewNotifyService(repo domain.NotificationRepository) *NotifyService {
	return &NotifyService{repo: repo}
}

// Enqueue 将通知持久化到发件箱表，状态为 "pending"。
// 调用者应在与业务事件相同的数据库事务中调用此方法，
// 以保证交付的原子性（确保在业务提交和发送之间崩溃时不会丢失事件）。
func (s *NotifyService) Enqueue(ctx context.Context, userID int64, channel, template, payload string) error {
	return s.repo.CreateOutbox(ctx, &domain.Notification{
		UserID:   userID,
		Channel:  channel,
		Template: template,
		Payload:  payload,
		Status:   NotificationStatusPending,
	})
}
