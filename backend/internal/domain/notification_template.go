package domain

import "time"

type NotificationChannel string

const (
	ChannelWechatSubscribe NotificationChannel = "wechat_subscribe"
	ChannelWechatTemplate  NotificationChannel = "wechat_template"
	ChannelSMS             NotificationChannel = "sms"
	ChannelInApp           NotificationChannel = "in_app"
)

type NotificationTemplate struct {
	ID        int64               `gorm:"primaryKey"`
	EventType string              `gorm:"size:50;index"` // order_created, order_paid, refund_success, travel_reminder
	Channel   NotificationChannel `gorm:"size:20"`
	Template  string              `gorm:"type:text"` // 模板内容，支持 {{.Field}} 占位符
	Enabled   bool                `gorm:"default:true"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

func (t *NotificationTemplate) Render(data map[string]string) string {
	result := t.Template
	for k, v := range data {
		placeholder := "{{." + k + "}}"
		result = replaceAll(result, placeholder, v)
	}
	return result
}

func replaceAll(s, old, new string) string {
	result := s
	for {
		if len(result) == 0 {
			break
		}
		i := find(result, old)
		if i == -1 {
			break
		}
		result = result[:i] + new + result[i+len(old):]
	}
	return result
}

func find(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
