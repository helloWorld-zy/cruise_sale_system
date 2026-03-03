package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// ImageRepository 提供图片资源的数据持久化实现。
type ImageRepository struct {
	db *gorm.DB // 数据库连接实例
}

// NewImageRepository 创建图片仓储实例。
func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// Create 新增图片记录。
func (r *ImageRepository) Create(ctx context.Context, img *domain.Image) error {
	return r.db.WithContext(ctx).Create(img).Error
}

// ListByEntity 按实体类型和实体 ID 查询图片，按排序和 ID 保持稳定顺序。
func (r *ImageRepository) ListByEntity(ctx context.Context, entityType string, entityID int64) ([]domain.Image, error) {
	var items []domain.Image
	err := r.db.WithContext(ctx).
		Model(&domain.Image{}).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("sort_order asc, id asc").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// DeleteByEntity 删除指定实体下的全部图片记录。
func (r *ImageRepository) DeleteByEntity(ctx context.Context, entityType string, entityID int64) error {
	return r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Delete(&domain.Image{}).Error
}

// UpdateSortOrder 更新单张图片的排序权重。
func (r *ImageRepository) UpdateSortOrder(ctx context.Context, id int64, sortOrder int) error {
	return r.db.WithContext(ctx).Model(&domain.Image{}).Where("id = ?", id).Update("sort_order", sortOrder).Error
}

// ReplaceImages 在事务内先删除旧图再批量插入新图，保证原子性。
func (r *ImageRepository) ReplaceImages(ctx context.Context, entityType string, entityID int64, images []*domain.Image) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity_type = ? AND entity_id = ?", entityType, entityID).Delete(&domain.Image{}).Error; err != nil {
			return err
		}
		for _, img := range images {
			if err := tx.Create(img).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
