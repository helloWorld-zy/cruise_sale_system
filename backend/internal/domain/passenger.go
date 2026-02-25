package domain

import "time"

// Passenger 表示用户可用于下单的出行乘客信息。
type Passenger struct {
	ID        int64  `gorm:"primaryKey"`
	UserID    int64  `gorm:"index"`
	Name      string `gorm:"size:50"`
	IDType    string `gorm:"size:20"`
	IDNumber  string `gorm:"size:50;column:id_number"`
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
