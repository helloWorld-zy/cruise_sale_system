package domain

import "time"

// Staff 表示系统后台管理人员（员工）信息。
// 员工通过用户名和密码登录管理后台，支持 RBAC 权限控制。
type Staff struct {
	ID           int64      `gorm:"primaryKey"`          // 主键 ID
	Username     string     `gorm:"size:50;uniqueIndex"` // 登录用户名（唯一）
	PasswordHash string     `gorm:"size:255"`            // 密码的 bcrypt 哈希值
	RealName     string     `gorm:"size:50"`             // 员工真实姓名
	Phone        string     `gorm:"size:20"`             // 联系电话
	Email        string     `gorm:"size:100"`            // 电子邮箱
	AvatarURL    string     `gorm:"size:500"`            // 头像图片地址
	Status       int16      `gorm:"default:1"`           // 状态：1=启用，0=停用
	LastLoginAt  *time.Time // 最后登录时间
	CreatedAt    time.Time  // 创建时间
	UpdatedAt    time.Time  // 更新时间
	DeletedAt    *time.Time `gorm:"index"` // 软删除时间
}
