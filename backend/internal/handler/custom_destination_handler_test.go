package handler

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type customDestinationHandlerRepoStub struct {
	listItems []domain.CustomDestination
	upserted  []domain.CustomDestination
}

func (s *customDestinationHandlerRepoStub) Create(_ context.Context, _ *domain.CustomDestination) error {
	return nil
}
func (s *customDestinationHandlerRepoStub) Update(_ context.Context, _ *domain.CustomDestination) error {
	return nil
}
func (s *customDestinationHandlerRepoStub) GetByID(_ context.Context, _ int64) (*domain.CustomDestination, error) {
	return nil, nil
}
func (s *customDestinationHandlerRepoStub) List(_ context.Context) ([]domain.CustomDestination, error) {
	return s.listItems, nil
}
func (s *customDestinationHandlerRepoStub) SearchByKeyword(_ context.Context, _ string) ([]domain.CustomDestination, error) {
	return nil, nil
}
func (s *customDestinationHandlerRepoStub) GetByLabel(_ context.Context, _, _ string) (*domain.CustomDestination, error) {
	return nil, nil
}
func (s *customDestinationHandlerRepoStub) Delete(_ context.Context, _ int64) error { return nil }
func (s *customDestinationHandlerRepoStub) UpsertByNameCountry(_ context.Context, dest *domain.CustomDestination) error {
	s.upserted = append(s.upserted, *dest)
	return nil
}

func TestCustomDestinationHandlerExportCSV(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &customDestinationHandlerRepoStub{listItems: []domain.CustomDestination{{Name: "迈阿密", Country: "美国", Latitude: float64PtrForHandler(25.7617), Longitude: float64PtrForHandler(-80.1918), Keywords: "迈阿密,miami", Status: 1, SortOrder: 100}}}
	h := NewCustomDestinationHandler(service.NewCustomDestinationService(repo))
	r := gin.New()
	r.GET("/custom-destinations/export", h.ExportCSV)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/custom-destinations/export", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Header().Get("Content-Type"), "text/csv") {
		t.Fatalf("expected csv content type, got %s", w.Header().Get("Content-Type"))
	}
	if !strings.Contains(w.Body.String(), "迈阿密") {
		t.Fatalf("expected csv body, got %s", w.Body.String())
	}
}

func TestCustomDestinationHandlerImportCSV(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &customDestinationHandlerRepoStub{}
	h := NewCustomDestinationHandler(service.NewCustomDestinationService(repo))
	r := gin.New()
	r.POST("/custom-destinations/import", h.ImportCSV)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "ports.csv")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("name,country,latitude,longitude,keywords,sort_order,status,description\n温哥华,加拿大,49.2827,-123.1207,\"温哥华,vancouver\",88,1,manual\n")); err != nil {
		t.Fatal(err)
	}
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/custom-destinations/import", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	if len(repo.upserted) != 1 {
		t.Fatalf("expected import to upsert one row, got %+v", repo.upserted)
	}
	if !strings.Contains(w.Body.String(), "imported") {
		t.Fatalf("expected import summary, got %s", w.Body.String())
	}
}

func float64PtrForHandler(v float64) *float64 { return &v }
