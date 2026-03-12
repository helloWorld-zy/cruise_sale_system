package service

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
)

type customDestinationCSVRepoStub struct {
	listItems []domain.CustomDestination
	upserted  []domain.CustomDestination
	upsertErr error
	listErr   error
}

func (s *customDestinationCSVRepoStub) Create(_ context.Context, _ *domain.CustomDestination) error {
	return nil
}
func (s *customDestinationCSVRepoStub) Update(_ context.Context, _ *domain.CustomDestination) error {
	return nil
}
func (s *customDestinationCSVRepoStub) GetByID(_ context.Context, _ int64) (*domain.CustomDestination, error) {
	return nil, nil
}
func (s *customDestinationCSVRepoStub) List(_ context.Context) ([]domain.CustomDestination, error) {
	return s.listItems, s.listErr
}
func (s *customDestinationCSVRepoStub) SearchByKeyword(_ context.Context, _ string) ([]domain.CustomDestination, error) {
	return nil, nil
}
func (s *customDestinationCSVRepoStub) GetByLabel(_ context.Context, _, _ string) (*domain.CustomDestination, error) {
	return nil, nil
}
func (s *customDestinationCSVRepoStub) Delete(_ context.Context, _ int64) error { return nil }
func (s *customDestinationCSVRepoStub) UpsertByNameCountry(_ context.Context, dest *domain.CustomDestination) error {
	if s.upsertErr != nil {
		return s.upsertErr
	}
	s.upserted = append(s.upserted, *dest)
	return nil
}

func TestCustomDestinationServiceExportCSV(t *testing.T) {
	repo := &customDestinationCSVRepoStub{listItems: []domain.CustomDestination{{
		Name: "迈阿密", Country: "美国", Latitude: float64Ptr(25.7617), Longitude: float64Ptr(-80.1918), Keywords: "迈阿密,miami", Description: "seed", SortOrder: 100, Status: 1,
	}}}
	svc := NewCustomDestinationService(repo)

	data, err := svc.ExportCSV(context.Background())
	if err != nil {
		t.Fatalf("ExportCSV returned error: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "name,country,latitude,longitude,keywords,sort_order,status,description") {
		t.Fatalf("expected csv header, got %s", content)
	}
	if !strings.Contains(content, "迈阿密,美国,25.7617,-80.1918,\"迈阿密,miami\",100,1,seed") {
		t.Fatalf("expected csv row, got %s", content)
	}
}

func TestCustomDestinationServiceImportCSV(t *testing.T) {
	repo := &customDestinationCSVRepoStub{}
	svc := NewCustomDestinationService(repo)
	input := bytes.NewBufferString("name,country,latitude,longitude,keywords,sort_order,status,description\n布宜诺斯艾利斯,阿根廷,-34.6037,-58.3816,\"布宜诺斯艾利斯,buenos aires\",90,1,manual\n")

	summary, err := svc.ImportCSV(context.Background(), input)
	if err != nil {
		t.Fatalf("ImportCSV returned error: %v", err)
	}
	if summary.Imported != 1 {
		t.Fatalf("expected 1 imported row, got %+v", summary)
	}
	if len(repo.upserted) != 1 {
		t.Fatalf("expected one upsert call, got %+v", repo.upserted)
	}
	if repo.upserted[0].Name != "布宜诺斯艾利斯" || repo.upserted[0].Country != "阿根廷" {
		t.Fatalf("expected imported row normalized, got %+v", repo.upserted[0])
	}
	if repo.upserted[0].Latitude == nil || *repo.upserted[0].Latitude != -34.6037 {
		t.Fatalf("expected latitude parsed, got %+v", repo.upserted[0])
	}
}

func TestCustomDestinationServiceImportCSVRejectsMissingCoords(t *testing.T) {
	repo := &customDestinationCSVRepoStub{}
	svc := NewCustomDestinationService(repo)
	input := bytes.NewBufferString("name,country,latitude,longitude,keywords,sort_order,status,description\n雷克雅未克,冰岛,,,reykjavik,80,1,manual\n")

	_, err := svc.ImportCSV(context.Background(), input)
	if err == nil || !strings.Contains(err.Error(), "row 2") {
		t.Fatalf("expected missing coordinate row error, got %v", err)
	}
}

func float64Ptr(v float64) *float64 { return &v }
