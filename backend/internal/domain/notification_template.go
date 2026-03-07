package domain

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
	"time"
)

// NotificationChannel 定义通知渠道类型。
type NotificationChannel string

const (
	ChannelWechatSubscribe NotificationChannel = "wechat_subscribe" // 微信关注通知
	ChannelWechatTemplate  NotificationChannel = "wechat_template"  // 微信模板消息
	ChannelSMS             NotificationChannel = "sms"              // 短信通知
	ChannelInApp           NotificationChannel = "in_app"           // 应用内通知
)

// NotificationTemplate 定义通知模板，支持事件类型和渲染变量。
// 模板使用 Go text/template 语法，支持 {{.Field}} 占位符。
type NotificationTemplate struct {
	ID        int64               `gorm:"primaryKey"`    // 主键 ID
	EventType string              `gorm:"size:50;index"` // 事件类型（如 order_created, order_paid, refund_success, travel_reminder）
	Channel   NotificationChannel `gorm:"size:20"`       // 通知渠道
	Template  string              `gorm:"type:text"`     // 模板内容，支持 {{.Field}} 占位符
	Enabled   bool                `gorm:"default:true"`  // 是否启用
	CreatedAt time.Time           `json:"created_at"`    // 创建时间
	UpdatedAt time.Time           `json:"updated_at"`    // 更新时间
}

var (
	// ErrNotificationTemplateInvalid 表示通知模板格式无效。
	ErrNotificationTemplateInvalid = errors.New("invalid notification template")
	// ErrNotificationTemplateVarNotAllowed 表示通知模板使用了未授权的变量。
	ErrNotificationTemplateVarNotAllowed = errors.New("notification template variable not allowed")
)

// defaultNotificationTemplateWhitelist 定义允许在模板中使用的变量白名单。
var defaultNotificationTemplateWhitelist = map[string]struct{}{
	"OrderNo":      {}, // 订单号
	"Name":         {}, // 姓名
	"Amount":       {}, // 金额
	"RefundAmount": {}, // 退款金额
	"VoyageName":   {}, // 航次名称
	"TravelDate":   {}, // 出行日期
	"Status":       {}, // 状态
}

// Render 使用给定数据渲染通知模板，返回最终的文本内容。
// 仅允许使用白名单中的变量，防止安全风险。
func (t *NotificationTemplate) Render(data map[string]string) (string, error) {
	tpl, err := template.New("notification").Option("missingkey=error").Parse(t.Template)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNotificationTemplateInvalid, err)
	}
	if err := validateTemplateWhitelist(tpl.Tree.Root, defaultNotificationTemplateWhitelist); err != nil {
		return "", err
	}
	b := strings.Builder{}
	if err := tpl.Execute(&b, data); err != nil {
		return "", fmt.Errorf("%w: %v", ErrNotificationTemplateInvalid, err)
	}
	return b.String(), nil
}

// validateTemplateWhitelist 递归验证模板 AST 节点，确保只使用白名单中的变量。
// 防止模板注入安全漏洞。
func validateTemplateWhitelist(node parse.Node, whitelist map[string]struct{}) error {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *parse.ListNode:
		for _, child := range n.Nodes {
			if err := validateTemplateWhitelist(child, whitelist); err != nil {
				return err
			}
		}
	case *parse.ActionNode:
		for _, cmd := range n.Pipe.Cmds {
			for _, arg := range cmd.Args {
				fieldNode, ok := arg.(*parse.FieldNode)
				if !ok || len(fieldNode.Ident) == 0 {
					continue
				}
				name := fieldNode.Ident[0]
				if _, ok := whitelist[name]; !ok {
					return fmt.Errorf("%w: %s", ErrNotificationTemplateVarNotAllowed, name)
				}
			}
		}
	case *parse.IfNode:
		if err := validateTemplateWhitelist(n.List, whitelist); err != nil {
			return err
		}
		return validateTemplateWhitelist(n.ElseList, whitelist)
	case *parse.RangeNode:
		if err := validateTemplateWhitelist(n.List, whitelist); err != nil {
			return err
		}
		return validateTemplateWhitelist(n.ElseList, whitelist)
	case *parse.WithNode:
		if err := validateTemplateWhitelist(n.List, whitelist); err != nil {
			return err
		}
		return validateTemplateWhitelist(n.ElseList, whitelist)
	}
	return nil
}
