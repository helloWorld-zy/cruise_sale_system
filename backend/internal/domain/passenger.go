package domain

import "time"

// Passenger 表示用户可用于下单的出行乘客信息。
// 一个用户可以维护多位乘客的证件资料，用于预订出行。
type Passenger struct {
	ID               int64     `gorm:"primaryKey"`                       // 主键 ID
	UserID           int64     `gorm:"index"`                            // 所属用户 ID
	Name             string    `gorm:"size:50"`                          // 乘客姓名（中文）
	EnglishName      string    `gorm:"size:100" json:"english_name"`     // 英文姓名
	IDType           string    `gorm:"size:20"`                          // 证件类型（身份证/护照等）
	IDNumber         string    `gorm:"size:50;column:id_number"`         // 证件号码
	Phone            string    `gorm:"size:20" json:"phone"`             // 手机号
	Email            string    `gorm:"size:100" json:"email"`            // 邮箱
	EmergencyContact string    `gorm:"size:50" json:"emergency_contact"` // 紧急联系人
	EmergencyPhone   string    `gorm:"size:20" json:"emergency_phone"`   // 紧急联系人电话
	SpecialNeeds     string    `gorm:"type:text" json:"special_needs"`   // 特殊需求备注
	Birthday         time.Time // 出生日期
	IsFavorite       bool      `gorm:"default:false" json:"is_favorite"` // 是否常用乘客
	CreatedAt        time.Time // 创建时间
	UpdatedAt        time.Time // 更新时间
}
