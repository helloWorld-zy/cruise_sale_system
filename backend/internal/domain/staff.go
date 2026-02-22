package domain

import "time"

type Staff struct {
	ID           int64  `gorm:"primaryKey"`
	Username     string `gorm:"size:50;uniqueIndex"`
	PasswordHash string `gorm:"size:255"`
	RealName     string `gorm:"size:50"`
	Phone        string `gorm:"size:20"`
	Email        string `gorm:"size:100"`
	AvatarURL    string `gorm:"size:500"`
	Status       int16  `gorm:"default:1"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}
