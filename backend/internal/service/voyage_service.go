package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

type VoyageRouteBuilder interface {
	BuildVoyageRouteMap(ctx context.Context, itineraries []domain.VoyageItinerary) (*domain.VoyageRouteMap, error)
}

type VoyageCityResolver interface {
	ResolveLabel(ctx context.Context, label string) (*ResolvedPortCity, error)
}

type VoyageService struct {
	repo         domain.VoyageRepository
	routeBuilder VoyageRouteBuilder
	cityResolver VoyageCityResolver
}

func NewVoyageService(repo domain.VoyageRepository, routeBuilder VoyageRouteBuilder) *VoyageService {
	return &VoyageService{repo: repo, routeBuilder: routeBuilder}
}

func (s *VoyageService) SetCityResolver(resolver VoyageCityResolver) *VoyageService {
	s.cityResolver = resolver
	return s
}

func (s *VoyageService) List(ctx context.Context) ([]domain.Voyage, error) {
	return s.repo.List(ctx)
}

func (s *VoyageService) ListPublic(ctx context.Context, cruiseID int64, keyword string, page, pageSize int) ([]domain.Voyage, int64, error) {
	return s.repo.ListPublic(ctx, cruiseID, keyword, page, pageSize)
}

func (s *VoyageService) Create(ctx context.Context, v *domain.Voyage) error {
	if err := s.enrichItineraryCoordinates(ctx, v); err != nil {
		return err
	}
	return s.repo.Create(ctx, v)
}

func (s *VoyageService) Update(ctx context.Context, v *domain.Voyage) error {
	if err := s.enrichItineraryCoordinates(ctx, v); err != nil {
		return err
	}
	return s.repo.Update(ctx, v)
}

func (s *VoyageService) GetByID(ctx context.Context, id int64) (*domain.Voyage, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil || len(item.Itineraries) == 0 {
		return item, nil
	}
	if s.routeBuilder != nil {
		routeMap, routeErr := s.routeBuilder.BuildVoyageRouteMap(ctx, item.Itineraries)
		if routeErr == nil {
			item.RouteMap = routeMap
		}
	}
	return item, nil
}

func (s *VoyageService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *VoyageService) enrichItineraryCoordinates(ctx context.Context, voyage *domain.Voyage) error {
	if s == nil || s.cityResolver == nil || voyage == nil {
		return nil
	}
	for index := range voyage.Itineraries {
		resolved, err := s.cityResolver.ResolveLabel(ctx, voyage.Itineraries[index].City)
		if err != nil {
			return err
		}
		if resolved == nil || resolved.IsSpecial {
			voyage.Itineraries[index].Latitude = nil
			voyage.Itineraries[index].Longitude = nil
			continue
		}
		voyage.Itineraries[index].Latitude = resolved.Latitude
		voyage.Itineraries[index].Longitude = resolved.Longitude
	}
	return nil
}
