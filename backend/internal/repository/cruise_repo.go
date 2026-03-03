package repository

import (
	"context"
	"fmt"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// CruiseRepository 提供邮轮实体的数据库操作。
type CruiseRepository struct {
	db *gorm.DB // 数据库连接实例
}

// NewCruiseRepository 创建邮轮仓储实例。
func NewCruiseRepository(db *gorm.DB) *CruiseRepository {
	return &CruiseRepository{db: db}
}

// Create 插入一条新的邮轮记录。
func (r *CruiseRepository) Create(ctx context.Context, cruise *domain.Cruise) error {
	return r.db.WithContext(ctx).Create(cruise).Error
}

// Update 保存邮轮的所有字段修改。
func (r *CruiseRepository) Update(ctx context.Context, cruise *domain.Cruise) error {
	return r.db.WithContext(ctx).Save(cruise).Error
}

// GetByID 根据主键查询邮轮记录。
func (r *CruiseRepository) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	var cruise domain.Cruise
	if err := r.db.WithContext(ctx).First(&cruise, id).Error; err != nil {
		return nil, err
	}
	return &cruise, nil
}

// List 分页查询邮轮列表，可按公司、关键词、状态筛选并支持排序。
func (r *CruiseRepository) List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	var items []domain.Cruise
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.Cruise{})
	if companyID > 0 {
		q = q.Where("company_id = ?", companyID)
	}
	if keyword != "" {
		q = q.Where("name LIKE ? OR english_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	// status 为 nil 表示不筛选；非 nil 时按值筛选（包括 0=下架）
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	orderBy := "sort_order desc, id desc"
	switch sortBy {
	case "tonnage_asc":
		orderBy = "tonnage asc, id desc"
	case "tonnage_desc":
		orderBy = "tonnage desc, id desc"
	case "name_asc":
		orderBy = "name asc, id desc"
	case "name_desc":
		orderBy = "name desc, id desc"
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order(orderBy).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Delete 软删除指定的邮轮记录。
func (r *CruiseRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Cruise{}, id).Error
}

// BatchUpdateStatus 批量更新邮轮状态，并在目标数量不匹配时回滚。
func (r *CruiseRepository) BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&domain.Cruise{}).Where("id IN ?", ids).Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if int(res.RowsAffected) != len(ids) {
			return fmt.Errorf("batch update cruise status affected=%d expected=%d", res.RowsAffected, len(ids))
		}
		return nil
	})
}
