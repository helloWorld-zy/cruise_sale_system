package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// CabinTypeMediaService 提供舱型媒体管理的业务能力。
type CabinTypeMediaService struct {
	repo domain.CabinTypeMediaRepository
}

func NewCabinTypeMediaService(repo domain.CabinTypeMediaRepository) *CabinTypeMediaService {
	return &CabinTypeMediaService{repo: repo}
}

func (s *CabinTypeMediaService) Create(ctx context.Context, media *domain.CabinTypeMedia) error {
	if media.IsPrimary {
		if err := s.repo.SetPrimary(ctx, media.CabinTypeID, media.MediaType, 0); err != nil {
			return err
		}
	}
	return s.repo.Create(ctx, media)
}

func (s *CabinTypeMediaService) Update(ctx context.Context, media *domain.CabinTypeMedia) error {
	if media.IsPrimary {
		if err := s.repo.SetPrimary(ctx, media.CabinTypeID, media.MediaType, media.ID); err != nil {
			return err
		}
	}
	return s.repo.Update(ctx, media)
}

func (s *CabinTypeMediaService) GetByID(ctx context.Context, id int64) (*domain.CabinTypeMedia, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CabinTypeMediaService) ListByCabinType(ctx context.Context, cabinTypeID int64) ([]domain.CabinTypeMedia, error) {
	return s.repo.ListByCabinType(ctx, cabinTypeID)
}

func (s *CabinTypeMediaService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *CabinTypeMediaService) SetPrimary(ctx context.Context, cabinTypeID int64, mediaType string, mediaID int64) error {
	return s.repo.SetPrimary(ctx, cabinTypeID, mediaType, mediaID)
}
