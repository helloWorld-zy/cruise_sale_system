package service

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubNotificationRepo 实现 domain.NotificationRepository。
type stubNotificationRepo struct {
	outbox []*domain.Notification
	sent   []int64
	failed []int64
}

func (r *stubNotificationRepo) CreateOutbox(_ context.Context, n *domain.Notification) error {
	r.outbox = append(r.outbox, n)
	return nil
}
func (r *stubNotificationRepo) ListPending(_ context.Context, _ int) ([]domain.Notification, error) {
	return nil, nil
}
func (r *stubNotificationRepo) MarkSent(_ context.Context, id int64) error {
	r.sent = append(r.sent, id)
	return nil
}
func (r *stubNotificationRepo) MarkFailed(_ context.Context, id int64) error {
	r.failed = append(r.failed, id)
	return nil
}

func TestNotifyService_Enqueue_SMS(t *testing.T) {
	repo := &stubNotificationRepo{}
	svc := NewNotifyService(repo)

	err := svc.Enqueue(context.Background(), 42, ChannelSMS, "booking_confirmed", `{"order_id":1}`)

	require.NoError(t, err)
	assert.Len(t, repo.outbox, 1)
	assert.Equal(t, int64(42), repo.outbox[0].UserID)
	assert.Equal(t, ChannelSMS, repo.outbox[0].Channel)
	assert.Equal(t, "booking_confirmed", repo.outbox[0].Template)
	assert.Equal(t, `{"order_id":1}`, repo.outbox[0].Payload)
	assert.Equal(t, NotificationStatusPending, repo.outbox[0].Status)
}

func TestNotifyService_Enqueue_WechatChannel(t *testing.T) {
	repo := &stubNotificationRepo{}
	svc := NewNotifyService(repo)

	err := svc.Enqueue(context.Background(), 1, ChannelWechat, "payment_success", `{}`)
	require.NoError(t, err)
	assert.Equal(t, ChannelWechat, repo.outbox[0].Channel)
}

func TestNotifyService_Enqueue_InboxChannel(t *testing.T) {
	repo := &stubNotificationRepo{}
	svc := NewNotifyService(repo)

	err := svc.Enqueue(context.Background(), 1, ChannelInbox, "order_update", `{}`)
	require.NoError(t, err)
	assert.Equal(t, ChannelInbox, repo.outbox[0].Channel)
}

func TestNotifyService_Enqueue_PersistsToOutbox(t *testing.T) {
	repo := &stubNotificationRepo{}
	svc := NewNotifyService(repo)

	// 多次入队应全部独立持久化。
	_ = svc.Enqueue(context.Background(), 1, ChannelSMS, "t1", `{}`)
	_ = svc.Enqueue(context.Background(), 2, ChannelWechat, "t2", `{}`)
	_ = svc.Enqueue(context.Background(), 3, ChannelInbox, "t3", `{}`)

	assert.Len(t, repo.outbox, 3)
}
