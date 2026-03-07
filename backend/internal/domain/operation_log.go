package domain

import "time"

// OperationLog 记录后台员工的系统操作日志，用于审计和追溯。
type OperationLog struct {
	ID         int64     `gorm:"primaryKey"`         // 主键 ID
	StaffID    int64     `gorm:"index"`              // 操作员工 ID
	Operation  string    `gorm:"size:50"`            // 操作类型（如 assign_role, create_booking）
	Resource   string    `gorm:"size:50"`            // 资源类型（如 staff, booking）
	ResourceID int64     `gorm:"column:resource_id"` // 资源 ID
	Details    string    `gorm:"type:text"`          // 操作详情（JSON 格式）
	IPAddress  string    `gorm:"size:50"`            // 客户端 IP 地址
	UserAgent  string    `gorm:"size:500"`           // 客户端 User-Agent
	CreatedAt  time.Time `gorm:"index"`              // 操作时间
}
