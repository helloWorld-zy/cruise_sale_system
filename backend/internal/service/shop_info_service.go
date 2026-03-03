package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

var (
	ErrInvalidShopInfoSingletonID = errors.New("shop info only supports singleton id=1")
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
	if info == nil {
		return nil
	}
	if info.ID != 0 && info.ID != 1 {
		return ErrInvalidShopInfoSingletonID
	}
	info.ID = 1
	return s.repo.Save(ctx, info)
}
