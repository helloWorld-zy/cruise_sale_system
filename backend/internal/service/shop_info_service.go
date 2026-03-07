package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

var (
	// ErrInvalidShopInfoSingletonID 表示商店信息仅支持 ID=1 的单例模式。
	ErrInvalidShopInfoSingletonID = errors.New("shop info only supports singleton id=1")
)

// ShopInfoRepository 定义商店信息数据访问接口。
type ShopInfoRepository interface {
	Get(ctx context.Context) (*domain.ShopInfo, error)     // 获取商店信息
	Save(ctx context.Context, info *domain.ShopInfo) error // 保存商店信息
}

// ShopInfoService 提供商店信息管理服务，支持单例模式的查询和更新。
type ShopInfoService struct {
	repo ShopInfoRepository // 商店信息数据仓储
}

// NewShopInfoService 创建商店信息服务实例。
func NewShopInfoService(repo ShopInfoRepository) *ShopInfoService {
	return &ShopInfoService{repo: repo}
}

// Get 获取商店信息（单例模式，始终返回 ID=1 的记录）。
func (s *ShopInfoService) Get(ctx context.Context) (*domain.ShopInfo, error) {
	return s.repo.Get(ctx)
}

// Update 更新商店信息。强制使用 ID=1，确保单例模式。
func (s *ShopInfoService) Update(ctx context.Context, info *domain.ShopInfo) error {
	if info == nil {
		return nil
	}
	if info.ID != 0 && info.ID != 1 {
		return ErrInvalidShopInfoSingletonID
	}
	info.ID = 1
	return s.repo.Save(ctx, info)
}
