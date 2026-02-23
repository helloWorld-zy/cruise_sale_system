package domain

import "time"

// Role 表示系统角色（如"admin"管理员、"operator"操作员等）。
// 用于 RBAC 权限控制，一个员工可以拥有多个角色。
type Role struct {
	ID          int64     `gorm:"primaryKey"`          // 主键 ID
	Name        string    `gorm:"size:50;uniqueIndex"` // 角色标识符（英文，全局唯一）
	DisplayName string    `gorm:"size:100"`            // 角色显示名称（中文）
	Description string    `gorm:"type:text"`           // 角色描述
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}
