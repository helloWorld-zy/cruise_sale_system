package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// fakeRouteService implements RouteService for testing.
type fakeRouteService struct{ routes []domain.Route }

func (f *fakeRouteService) List(_ context.Context) ([]domain.Route, error) { return f.routes, nil }
func (f *fakeRouteService) Create(_ context.Context, r *domain.Route) error {
	f.routes = append(f.routes, *r)
	return nil
}
func (f *fakeRouteService) Update(_ context.Context, r *domain.Route) error { return nil }
func (f *fakeRouteService) GetByID(_ context.Context, id int64) (*domain.Route, error) {
	return &domain.Route{}, nil
}
func (f *fakeRouteService) Delete(_ context.Context, id int64) error { return nil }

func TestRouteListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := &fakeRouteService{routes: []domain.Route{{ID: 1, Code: "R1", Name: "Route 1"}}}
	h := NewRouteHandler(svc)
	r.GET("/admin/routes", h.List)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/routes", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRouteDeleteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewRouteHandler(&fakeRouteService{})
	r.DELETE("/admin/routes/:id", h.Delete)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/admin/routes/1", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

// fakeVoyageService implements VoyageService for testing.
type fakeVoyageService struct{}

func (f *fakeVoyageService) ListByRoute(_ context.Context, routeID int64) ([]domain.Voyage, error) {
	return []domain.Voyage{}, nil
}
func (f *fakeVoyageService) Create(_ context.Context, v *domain.Voyage) error { return nil }
func (f *fakeVoyageService) Update(_ context.Context, v *domain.Voyage) error { return nil }
func (f *fakeVoyageService) Delete(_ context.Context, id int64) error         { return nil }

func TestVoyageListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewVoyageHandler(&fakeVoyageService{})
	r.GET("/admin/voyages", h.List)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/voyages", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
