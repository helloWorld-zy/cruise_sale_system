package repository

import (
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// UserRepository 提供 C 端用户的数据库操作。
type UserRepository struct{ db *gorm.DB }

// NewUserRepository 创建用户仓储实例。
func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

// FindOrCreateByPhone 根据手机号查找用户；若不存在则创建并返回。
// 保证幂等：同一手机号多次调用只创建一条记录。
func (r *UserRepository) FindOrCreateByPhone(phone string) (*domain.User, error) {
	var u domain.User
	result := r.db.Where("phone = ?", phone).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		u = domain.User{Phone: phone, Status: 1}
		if err := r.db.Create(&u).Error; err != nil {
			// 并发场景下可能被其他请求先创建，重查一次保证幂等。
			if err2 := r.db.Where("phone = ?", phone).First(&u).Error; err2 != nil {
				return nil, err
			}
		}
		return &u, nil
	}
	return &u, result.Error
}
