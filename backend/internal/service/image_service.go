package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

// ErrMultiplePrimaryImages 当图片列表中有多张主图时返回此错误。
var ErrMultiplePrimaryImages = errors.New("only one primary image allowed per entity")

// ImageInput 描述图片画廊写入参数。
type ImageInput struct {
	URL       string
	SortOrder int
	IsPrimary bool
}

// ImageService 负责多实体图片画廊管理。
type ImageService struct {
	repo domain.ImageRepository // 图片仓储
}

// NewImageService 创建图片服务实例。
func NewImageService(repo domain.ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

// SetImages 覆盖式设置实体图片：在事务内先清空旧图，再按输入顺序重建。
// 校验：同一实体下最多只能有一张主图（IsPrimary=true）。
func (s *ImageService) SetImages(ctx context.Context, entityType string, entityID int64, images []ImageInput) error {
	// 校验主图唯一性
	primaryCount := 0
	for _, in := range images {
		if in.IsPrimary {
			primaryCount++
		}
	}
	if primaryCount > 1 {
		return ErrMultiplePrimaryImages
	}

	// 构建领域对象列表
	domainImages := make([]*domain.Image, 0, len(images))
	for _, in := range images {
		domainImages = append(domainImages, &domain.Image{
			EntityType: entityType,
			EntityID:   entityID,
			URL:        in.URL,
			SortOrder:  in.SortOrder,
			IsPrimary:  in.IsPrimary,
		})
	}

	return s.repo.ReplaceImages(ctx, entityType, entityID, domainImages)
}

// ListImages 查询实体关联的全部图片。
func (s *ImageService) ListImages(ctx context.Context, entityType string, entityID int64) ([]domain.Image, error) {
	return s.repo.ListByEntity(ctx, entityType, entityID)
}
