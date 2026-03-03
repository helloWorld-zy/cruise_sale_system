package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// PassengerRepo 定义乘客数据访问接口。
type PassengerRepo interface {
	ListByUser(ctx context.Context, userID int64) ([]domain.Passenger, error)
	UpdateFavorite(ctx context.Context, id int64, isFavorite bool) error
}

// PassengerService 提供常用乘客管理功能。
type PassengerService struct{ repo PassengerRepo }

// NewPassengerService 创建乘客服务实例。
func NewPassengerService(repo PassengerRepo) *PassengerService { return &PassengerService{repo: repo} }

// ListFavorites 查询用户的所有常用乘客。
func (s *PassengerService) ListFavorites(ctx context.Context, userID int64) ([]domain.Passenger, error) {
	all, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var favorites []domain.Passenger
	for _, p := range all {
		if p.IsFavorite {
			favorites = append(favorites, p)
		}
	}
	return favorites, nil
}

// ToggleFavorite 切换乘客的常用状态。
func (s *PassengerService) ToggleFavorite(ctx context.Context, id int64, isFavorite bool) error {
	return s.repo.UpdateFavorite(ctx, id, isFavorite)
}
