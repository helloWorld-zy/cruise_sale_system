package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

type customDestinationRepoStub struct {
	searchResults []domain.CustomDestination
}

func (s *customDestinationRepoStub) Create(_ context.Context, _ *domain.CustomDestination) error {
	return nil
}
func (s *customDestinationRepoStub) Update(_ context.Context, _ *domain.CustomDestination) error {
	return nil
}
func (s *customDestinationRepoStub) GetByID(_ context.Context, _ int64) (*domain.CustomDestination, error) {
	return nil, nil
}
func (s *customDestinationRepoStub) List(_ context.Context) ([]domain.CustomDestination, error) {
	return nil, nil
}
func (s *customDestinationRepoStub) SearchByKeyword(_ context.Context, keyword string) ([]domain.CustomDestination, error) {
	trimmed := strings.TrimSpace(keyword)
	if trimmed == "" {
		return nil, nil
	}
	items := make([]domain.CustomDestination, 0, len(s.searchResults))
	for _, item := range s.searchResults {
		if strings.Contains(item.Name, trimmed) || strings.Contains(item.Country, trimmed) || strings.Contains(strings.ToLower(item.Keywords), strings.ToLower(trimmed)) {
			items = append(items, item)
		}
	}
	return items, nil
}
func (s *customDestinationRepoStub) GetByLabel(_ context.Context, name, country string) (*domain.CustomDestination, error) {
	for _, item := range s.searchResults {
		if item.Name == name && item.Country == country {
			copyItem := item
			return &copyItem, nil
		}
	}
	return nil, nil
}
func (s *customDestinationRepoStub) UpsertByNameCountry(_ context.Context, _ *domain.CustomDestination) error {
	return nil
}
func (s *customDestinationRepoStub) Delete(_ context.Context, _ int64) error { return nil }

func TestPortCityServiceSearchFormatsCityCountryAndIncludesSeaCruise(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"lat": "37.4563",
				"lon": "126.7052",
				"namedetails": map[string]any{
					"name:zh-Hans": "仁川",
				},
				"address": map[string]any{
					"city":         "仁川;仁川廣域市",
					"country":      "韩国;韓國",
					"country_code": "kr",
				},
			},
		})
	}))
	defer server.Close()

	svc := NewPortCityService(PortCityServiceConfig{Endpoint: server.URL, Timeout: 2 * time.Second})
	items, err := svc.Search(context.Background(), "仁川")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(items) == 0 || items[0].Label != "仁川（韩国）" {
		t.Fatalf("expected formatted city-country label, got %+v", items)
	}

	special, err := svc.Search(context.Background(), "海上")
	if err != nil {
		t.Fatalf("Search special returned error: %v", err)
	}
	if len(special) == 0 || special[0].Label != "海上巡游" || !special[0].IsSpecial {
		t.Fatalf("expected sea cruise special option, got %+v", special)
	}
}

func TestPortCityServiceSearchUsesSimplifiedChineseAndSingleCharLocalMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"lat":  "25.7741",
				"lon":  "-80.1935",
				"name": "迈阿密;邁阿密",
				"namedetails": map[string]any{
					"name:zh":      "迈阿密;邁阿密",
					"name:zh-Hans": "迈阿密",
				},
				"address": map[string]any{
					"city":         "迈阿密;邁阿密",
					"country":      "美国;美國",
					"country_code": "us",
				},
			},
			{
				"lat":  "20.4320",
				"lon":  "-86.9206",
				"name": "Isla Cozumel",
				"address": map[string]any{
					"island":       "Isla Cozumel",
					"country":      "墨西哥",
					"country_code": "mx",
				},
			},
		})
	}))
	defer server.Close()

	svc := NewPortCityService(PortCityServiceConfig{Endpoint: server.URL, Timeout: 2 * time.Second})
	svc.SetCustomDestinationRepo(&customDestinationRepoStub{searchResults: []domain.CustomDestination{{Name: "迈阿密", Country: "美国", Latitude: floatPtr(25.7617), Longitude: floatPtr(-80.1918), Keywords: "迈阿密,miami"}}})
	items, err := svc.Search(context.Background(), "迈阿密")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(items) == 0 || items[0].Label != "迈阿密（美国）" {
		t.Fatalf("expected simplified chinese label, got %+v", items)
	}
	for _, item := range items {
		if item.Label == "Isla Cozumel（墨西哥）" {
			t.Fatalf("expected english-only remote result to be filtered, got %+v", items)
		}
	}

	cozumelItems, err := svc.Search(context.Background(), "科")
	if err != nil {
		t.Fatalf("Search single-char local returned error: %v", err)
	}
	if len(cozumelItems) != 0 {
		t.Fatalf("expected no built-in local catalog match without dictionary repo seed, got %+v", cozumelItems)
	}

	svc.SetCustomDestinationRepo(&customDestinationRepoStub{searchResults: []domain.CustomDestination{{Name: "科苏梅尔", Country: "墨西哥", Latitude: floatPtr(20.4229839), Longitude: floatPtr(-86.9223432), Keywords: "科苏梅尔,cozumel"}}})
	cozumelItems, err = svc.Search(context.Background(), "科")
	if err != nil {
		t.Fatalf("Search single-char dictionary returned error: %v", err)
	}
	if len(cozumelItems) == 0 || cozumelItems[0].Label != "科苏梅尔（墨西哥）" {
		t.Fatalf("expected single-char local match, got %+v", cozumelItems)
	}
}

func floatPtr(v float64) *float64 { return &v }

func TestPortCityServiceResolveLabelReturnsCoordinates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"lat": "37.4563",
				"lon": "126.7052",
				"address": map[string]any{
					"city":    "仁川",
					"country": "韩国",
				},
			},
		})
	}))
	defer server.Close()

	svc := NewPortCityService(PortCityServiceConfig{Endpoint: server.URL, Timeout: 2 * time.Second})
	resolved, err := svc.ResolveLabel(context.Background(), "仁川（韩国）")
	if err != nil {
		t.Fatalf("ResolveLabel returned error: %v", err)
	}
	if resolved == nil || resolved.Latitude == nil || resolved.Longitude == nil || *resolved.Latitude != 37.4563 || *resolved.Longitude != 126.7052 {
		t.Fatalf("expected resolved coordinates, got %+v", resolved)
	}

	special, err := svc.ResolveLabel(context.Background(), "海上巡游")
	if err != nil {
		t.Fatalf("ResolveLabel special returned error: %v", err)
	}
	if special == nil || !special.IsSpecial || special.Latitude != nil || special.Longitude != nil {
		t.Fatalf("expected special option without coordinates, got %+v", special)
	}
}
