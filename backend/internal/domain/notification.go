package domain

import "time"

// Notification 表示发件箱通知队列中的一条消息。
// 遵循发件箱模式：在业务事务中将通知写入此表，
// 然后由后台调度程序异步传递。
type Notification struct {
	ID        int64     `gorm:"primaryKey"` // 主键 ID
	UserID    int64     `gorm:"index"`      // 接收通知的用户 ID
	Channel   string    `gorm:"size:20"`    // 通知渠道（sms / wechat / inbox）
	Template  string    `gorm:"size:50"`    // 通知模板标识符
	Payload   string    `gorm:"type:text"`  // 通知负载（JSON 格式）
	Status    string    `gorm:"size:20"`    // 状态（pending / sent / failed）
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}

// TableName 显式设置表名，以免 GORM 出现意外的复数形式。
func (Notification) TableName() string { return "notifications" }
