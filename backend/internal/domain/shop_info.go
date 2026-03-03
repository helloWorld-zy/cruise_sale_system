package domain

type ShopInfo struct {
	ID              int64  `gorm:"primaryKey"`
	Name            string `gorm:"size:100"`
	Logo            string `gorm:"size:500"`
	ContactPhone    string `gorm:"size:20"`
	ContactEmail    string `gorm:"size:100"`
	CompanyDesc     string `gorm:"type:text"`
	ServiceDesc     string `gorm:"type:text"`
	ICPNumber       string `gorm:"size:50"`
	BusinessLicense string `gorm:"size:100"`
	Address         string `gorm:"size:200"`
	Wechat          string `gorm:"size:50"`
}
