package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// StaffRepository 提供员工实体的数据库操作。
type StaffRepository struct {
	db *gorm.DB // 数据库连接实例
}

// NewStaffRepository 创建员工仓储实例。
func NewStaffRepository(db *gorm.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

// Create 插入一条新的员工记录。
func (r *StaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

// GetByUsername 根据用户名查询员工（用于登录验证），排除已软删除的记录。
func (r *StaffRepository) GetByUsername(ctx context.Context, username string) (*domain.Staff, error) {
	var staff domain.Staff
	if err := r.db.WithContext(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// GetByID 根据主键查询员工记录。
func (r *StaffRepository) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	var staff domain.Staff
	if err := r.db.WithContext(ctx).First(&staff, id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// Update 保存员工的所有字段修改。
func (r *StaffRepository) Update(ctx context.Context, staff *domain.Staff) error {
	return r.db.WithContext(ctx).Save(staff).Error
}

// Delete 软删除指定的员工记录。
func (r *StaffRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Staff{}, id).Error
}
