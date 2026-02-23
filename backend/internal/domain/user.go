package domain

import "time"

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Phone     string `gorm:"size:20;uniqueIndex"`
	WxOpenID  string `gorm:"size:80;uniqueIndex"`
	Nickname  string `gorm:"size:50"`
	AvatarURL string `gorm:"size:500"`
	Status    int16  `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
