package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ShopInfoRepository 提供店铺信息单行配置读写。
type ShopInfoRepository struct {
	db *gorm.DB
}

func NewShopInfoRepository(db *gorm.DB) *ShopInfoRepository {
	return &ShopInfoRepository{db: db}
}

func (r *ShopInfoRepository) Get(ctx context.Context) (*domain.ShopInfo, error) {
	var info domain.ShopInfo
	err := r.db.WithContext(ctx).First(&info, 1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *ShopInfoRepository) Save(ctx context.Context, info *domain.ShopInfo) error {
	if info == nil {
		return nil
	}
	// 单行配置：固定主键 1，重复写入走更新。
	info.ID = 1
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "logo", "contact_phone", "contact_email", "company_desc", "service_desc", "icp_number", "business_license", "address", "wechat"}),
	}).Create(info).Error
}
