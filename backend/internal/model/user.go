package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	BaseModel
	Phone        string `gorm:"uniqueIndex;not null" json:"phone"`
	WxOpenID     string `gorm:"index" json:"wx_openid"`
	PasswordHash string `json:"-"`
}

type Staff struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	RoleID   string `json:"role_id"` // Casbin role identifier
}
