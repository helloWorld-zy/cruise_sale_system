package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNotification_TableName 验证通知模型表名映射。
func TestNotification_TableName(t *testing.T) {
	n := Notification{}
	assert.Equal(t, "notifications", n.TableName())
}
