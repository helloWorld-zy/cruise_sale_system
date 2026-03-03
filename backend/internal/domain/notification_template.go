package domain

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
	"time"
)

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

var (
	ErrNotificationTemplateInvalid       = errors.New("invalid notification template")
	ErrNotificationTemplateVarNotAllowed = errors.New("notification template variable not allowed")
)

var defaultNotificationTemplateWhitelist = map[string]struct{}{
	"OrderNo":      {},
	"Name":         {},
	"Amount":       {},
	"RefundAmount": {},
	"VoyageName":   {},
	"TravelDate":   {},
	"Status":       {},
}

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
