package domain

import "time"

// User 表示 C 端登录用户基础资料。
// 支持手机号、微信 OpenID 和支付宝 UID 三种唯一标识，均可用于登录。
type User struct {
	ID        int64     `gorm:"primaryKey"`                            // 主键 ID
	Phone     string    `gorm:"size:20;uniqueIndex"`                   // 手机号（唯一）
	WxOpenID  string    `gorm:"size:80;uniqueIndex"`                   // 微信 OpenID（唯一）
	AlipayUID string    `gorm:"size:80;uniqueIndex" json:"alipay_uid"` // 支付宝用户ID（唯一）
	Email     string    `gorm:"size:100" json:"email"`                 // 邮箱
	Nickname  string    `gorm:"size:50"`                               // 用户昵称
	AvatarURL string    `gorm:"size:500"`                              // 头像图片地址
	Status    int16     `gorm:"default:1"`                             // 状态：1=启用，0=停用
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}
