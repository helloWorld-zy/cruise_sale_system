package domain

// ShopInfo 表示商城基础信息（单例模式，ID 固定为 1）。
// 包含公司名称、联系方式、ICP 备案号、营业执照等基本信息。
type ShopInfo struct {
	ID              int64  `gorm:"primaryKey"` // 主键 ID（固定为 1）
	Name            string `gorm:"size:100"`   // 商城名称
	Logo            string `gorm:"size:500"`   // Logo 图片地址
	ContactPhone    string `gorm:"size:20"`    // 联系电话
	ContactEmail    string `gorm:"size:100"`   // 联系邮箱
	CompanyDesc     string `gorm:"type:text"`  // 公司简介
	ServiceDesc     string `gorm:"type:text"`  // 服务说明
	ICPNumber       string `gorm:"size:50"`    // ICP 备案号
	BusinessLicense string `gorm:"size:100"`   // 营业执照号
	Address         string `gorm:"size:200"`   // 公司地址
	Wechat          string `gorm:"size:50"`    // 微信公众号
}
