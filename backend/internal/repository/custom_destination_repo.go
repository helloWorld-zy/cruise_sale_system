package repository

import (
	"context"
	"strings"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CustomDestinationRepository 提供自定义目的地的数据库操作。
type CustomDestinationRepository struct {
	db *gorm.DB
}

// NewCustomDestinationRepository 创建自定义目的地仓储实例。
func NewCustomDestinationRepository(db *gorm.DB) *CustomDestinationRepository {
	return &CustomDestinationRepository{db: db}
}

// Create 插入一条新的自定义目的地记录。
func (r *CustomDestinationRepository) Create(ctx context.Context, dest *domain.CustomDestination) error {
	return r.db.WithContext(ctx).Create(dest).Error
}

// Update 保存自定义目的地的所有字段修改。
func (r *CustomDestinationRepository) Update(ctx context.Context, dest *domain.CustomDestination) error {
	return r.db.WithContext(ctx).Save(dest).Error
}

// GetByID 根据主键查询自定义目的地。
func (r *CustomDestinationRepository) GetByID(ctx context.Context, id int64) (*domain.CustomDestination, error) {
	var item domain.CustomDestination
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// List 查询所有自定义目的地，按排序权重和 ID 降序排列。
func (r *CustomDestinationRepository) List(ctx context.Context) ([]domain.CustomDestination, error) {
	var items []domain.CustomDestination
	if err := r.db.WithContext(ctx).Order("sort_order desc, id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// SearchByKeyword 按关键词搜索启用状态的自定义目的地（匹配 name、country、keywords 字段）。
func (r *CustomDestinationRepository) SearchByKeyword(ctx context.Context, keyword string) ([]domain.CustomDestination, error) {
	trimmed := strings.TrimSpace(keyword)
	if trimmed == "" {
		return nil, nil
	}
	var items []domain.CustomDestination
	pattern := "%" + trimmed + "%"
	err := r.db.WithContext(ctx).
		Where("status = 1").
		Where("name ILIKE ? OR country ILIKE ? OR keywords ILIKE ?", pattern, pattern, pattern).
		Order("sort_order desc, id desc").
		Limit(10).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetByLabel 根据 "名称（国家）" 格式的标签查找启用状态的自定义目的地。
func (r *CustomDestinationRepository) GetByLabel(ctx context.Context, name, country string) (*domain.CustomDestination, error) {
	var item domain.CustomDestination
	err := r.db.WithContext(ctx).
		Where("status = 1 AND name = ? AND country = ?", name, country).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpsertByNameCountry 按名称+国家更新或新增自定义目的地。
func (r *CustomDestinationRepository) UpsertByNameCountry(ctx context.Context, dest *domain.CustomDestination) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "name"}, {Name: "country"}},
			DoUpdates: clause.Assignments(map[string]any{
				"latitude":    dest.Latitude,
				"longitude":   dest.Longitude,
				"keywords":    dest.Keywords,
				"description": dest.Description,
				"status":      dest.Status,
				"sort_order":  dest.SortOrder,
				"updated_at":  gorm.Expr("NOW()"),
				"deleted_at":  nil,
			}),
		}).
		Create(dest).Error
}

// Delete 软删除指定的自定义目的地。
func (r *CustomDestinationRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.CustomDestination{}, id).Error
}
