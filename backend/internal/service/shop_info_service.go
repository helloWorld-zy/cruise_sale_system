package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

type ShopInfoRepository interface {
	Get(ctx context.Context) (*domain.ShopInfo, error)
	Save(ctx context.Context, info *domain.ShopInfo) error
}

type ShopInfoService struct {
	repo ShopInfoRepository
}

func NewShopInfoService(repo ShopInfoRepository) *ShopInfoService {
	return &ShopInfoService{repo: repo}
}

func (s *ShopInfoService) Get(ctx context.Context) (*domain.ShopInfo, error) {
	return s.repo.Get(ctx)
}

func (s *ShopInfoService) Update(ctx context.Context, info *domain.ShopInfo) error {
	return s.repo.Save(ctx, info)
}
