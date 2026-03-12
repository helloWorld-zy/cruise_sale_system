package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type voyageRepoStub struct {
	item    *domain.Voyage
	created *domain.Voyage
	updated *domain.Voyage
}

func (s *voyageRepoStub) Create(_ context.Context, item *domain.Voyage) error {
	s.created = item
	return nil
}
func (s *voyageRepoStub) Update(_ context.Context, item *domain.Voyage) error {
	s.updated = item
	return nil
}
func (s *voyageRepoStub) Delete(context.Context, int64) error { return nil }
func (s *voyageRepoStub) List(context.Context) ([]domain.Voyage, error) {
	return []domain.Voyage{}, nil
}
func (s *voyageRepoStub) ListPublic(context.Context, int64, string, int, int) ([]domain.Voyage, int64, error) {
	return []domain.Voyage{}, 0, nil
}
func (s *voyageRepoStub) GetByID(context.Context, int64) (*domain.Voyage, error) {
	return s.item, nil
}

type maritimeRouteBuilderStub struct {
	routeMap *domain.VoyageRouteMap
}

func (s *maritimeRouteBuilderStub) BuildVoyageRouteMap(context.Context, []domain.VoyageItinerary) (*domain.VoyageRouteMap, error) {
	return s.routeMap, nil
}

type voyageCityResolverStub struct{}

func (s *voyageCityResolverStub) ResolveLabel(_ context.Context, label string) (*ResolvedPortCity, error) {
	switch label {
	case "仁川（韩国）":
		lat, lon := 37.4563, 126.7052
		return &ResolvedPortCity{Label: label, Latitude: &lat, Longitude: &lon}, nil
	case "海上巡游":
		return &ResolvedPortCity{Label: label, IsSpecial: true}, nil
	default:
		return nil, nil
	}
}

func TestSeaRouteClientBuildVoyageRouteMapUsesExplicitCoordinates(t *testing.T) {
	requested := struct {
		opos string
		dpos string
		res  string
	}{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested.opos = r.URL.Query().Get("opos")
		requested.dpos = r.URL.Query().Get("dpos")
		requested.res = r.URL.Query().Get("res")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"dist":   812.4,
			"geom": map[string]any{
				"type": "MultiLineString",
				"coordinates": [][][]float64{{
					{121.4737, 31.2304},
					{124.5000, 32.1000},
					{130.4017, 33.5902},
				}},
			},
		})
	}))
	defer server.Close()

	client := NewSeaRouteClient(SeaRouteClientConfig{Endpoint: server.URL, Timeout: 2 * time.Second, ResolutionKM: 20})
	lat1, lon1 := 31.2304, 121.4737
	lat2, lon2 := 33.5902, 130.4017
	routeMap, err := client.BuildVoyageRouteMap(context.Background(), []domain.VoyageItinerary{
		{DayNo: 1, StopIndex: 1, City: "任意起点", Latitude: &lat1, Longitude: &lon1},
		{DayNo: 2, StopIndex: 1, City: "海上巡游"},
		{DayNo: 3, StopIndex: 1, City: "任意终点", Latitude: &lat2, Longitude: &lon2},
	})
	if err != nil {
		t.Fatalf("BuildVoyageRouteMap returned error: %v", err)
	}

	if requested.opos != "121.473700,31.230400" {
		t.Fatalf("unexpected opos: %s", requested.opos)
	}
	if requested.dpos != "130.401700,33.590200" {
		t.Fatalf("unexpected dpos: %s", requested.dpos)
	}
	if requested.res != "20" {
		t.Fatalf("unexpected res: %s", requested.res)
	}
	if routeMap == nil || len(routeMap.Coordinates) != 1 {
		t.Fatalf("expected normalized route map, got %+v", routeMap)
	}
	if routeMap.DistanceKM != 812.4 {
		t.Fatalf("expected distance 812.4, got %+v", routeMap)
	}
}

func TestVoyageServiceGetByIDEnrichesRouteMap(t *testing.T) {
	repo := &voyageRepoStub{item: &domain.Voyage{ID: 7, Code: "VOY-7", Itineraries: []domain.VoyageItinerary{{DayNo: 1, StopIndex: 1, City: "上海"}, {DayNo: 2, StopIndex: 1, City: "福冈"}}}}
	svc := NewVoyageService(repo, &maritimeRouteBuilderStub{routeMap: &domain.VoyageRouteMap{Provider: "searoute", GeometryType: "MultiLineString", Coordinates: [][][]float64{{{121.4, 31.2}, {130.4, 33.5}}}}})

	item, err := svc.GetByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if item.RouteMap == nil {
		t.Fatalf("expected route_map enrichment, got %+v", item)
	}
	if item.RouteMap.Provider != "searoute" {
		t.Fatalf("expected searoute provider, got %+v", item.RouteMap)
	}
}

func TestVoyageServiceCreateResolvesSelectedCityLabelsToCoordinates(t *testing.T) {
	repo := &voyageRepoStub{}
	svc := NewVoyageService(repo, nil)
	svc.SetCityResolver(&voyageCityResolverStub{})

	err := svc.Create(context.Background(), &domain.Voyage{Itineraries: []domain.VoyageItinerary{
		{DayNo: 1, StopIndex: 1, City: "仁川（韩国）"},
		{DayNo: 2, StopIndex: 1, City: "海上巡游"},
	}})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if repo.created == nil || repo.created.Itineraries[0].Latitude == nil || repo.created.Itineraries[0].Longitude == nil {
		t.Fatalf("expected first itinerary coordinates to be filled, got %+v", repo.created)
	}
	if repo.created.Itineraries[1].Latitude != nil || repo.created.Itineraries[1].Longitude != nil {
		t.Fatalf("expected sea cruise stop to keep empty coordinates, got %+v", repo.created.Itineraries[1])
	}
}
