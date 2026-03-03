package service

import (
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNotifyTemplateRender(t *testing.T) {
	tpl := domain.NotificationTemplate{
		EventType: "order_paid",
		Channel:   domain.ChannelSMS,
		Template:  "您的订单{{.OrderNo}}已支付成功",
	}
	result := tpl.Render(map[string]string{"OrderNo": "ORD001"})
	if result != "您的订单ORD001已支付成功" {
		t.Fatal("template render failed")
	}
	assert.Equal(t, "您的订单ORD001已支付成功", result)
}

func TestNotifyTemplateRender_Multiple(t *testing.T) {
	tpl := domain.NotificationTemplate{
		EventType: "order_paid",
		Channel:   domain.ChannelSMS,
		Template:  "{{.Name}}您好，订单{{.OrderNo}}已支付{{.Amount}}元",
	}
	result := tpl.Render(map[string]string{"Name": "张三", "OrderNo": "ORD001", "Amount": "100"})
	assert.Equal(t, "张三您好，订单ORD001已支付100元", result)
}
