package domain

import "time"

type OperationLog struct {
	ID         int64  `gorm:"primaryKey"`
	StaffID    int64  `gorm:"index"`
	Operation  string `gorm:"size:50"`
	Resource   string `gorm:"size:50"`
	ResourceID int64
	Details    string    `gorm:"type:text"`
	IPAddress  string    `gorm:"size:50"`
	UserAgent  string    `gorm:"size:500"`
	CreatedAt  time.Time `gorm:"index"`
}
