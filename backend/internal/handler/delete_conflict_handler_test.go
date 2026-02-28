package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

const fkErrMsg = "pq: update or delete on table \"x\" violates foreign key constraint \"x_ref_fkey\" on table \"y\""

type routeDeleteErrSvc struct{ err error }

func (s *routeDeleteErrSvc) List(context.Context) ([]domain.Route, error)                { return nil, nil }
func (s *routeDeleteErrSvc) Create(context.Context, *domain.Route) error                 { return nil }
func (s *routeDeleteErrSvc) Update(context.Context, *domain.Route) error                 { return nil }
func (s *routeDeleteErrSvc) GetByID(context.Context, int64) (*domain.Route, error)       { return &domain.Route{}, nil }
func (s *routeDeleteErrSvc) Delete(context.Context, int64) error                          { return s.err }

type voyageDeleteErrSvc struct{ err error }

func (s *voyageDeleteErrSvc) ListByRoute(context.Context, int64) ([]domain.Voyage, error) { return nil, nil }
func (s *voyageDeleteErrSvc) Create(context.Context, *domain.Voyage) error                 { return nil }
func (s *voyageDeleteErrSvc) Update(context.Context, *domain.Voyage) error                 { return nil }
func (s *voyageDeleteErrSvc) GetByID(context.Context, int64) (*domain.Voyage, error)      { return &domain.Voyage{}, nil }
func (s *voyageDeleteErrSvc) Delete(context.Context, int64) error                          { return s.err }

type cabinDeleteErrSvc struct{ err error }

func (s *cabinDeleteErrSvc) ListByVoyage(context.Context, int64) ([]domain.CabinSKU, error) { return nil, nil }
func (s *cabinDeleteErrSvc) GetByID(context.Context, int64) (*domain.CabinSKU, error)       { return &domain.CabinSKU{}, nil }
func (s *cabinDeleteErrSvc) Create(context.Context, *domain.CabinSKU) error                  { return nil }
func (s *cabinDeleteErrSvc) Update(context.Context, *domain.CabinSKU) error                  { return nil }
func (s *cabinDeleteErrSvc) Delete(context.Context, int64) error                             { return s.err }
func (s *cabinDeleteErrSvc) GetInventory(context.Context, int64) (domain.CabinInventory, error) {
	return domain.CabinInventory{}, nil
}
func (s *cabinDeleteErrSvc) AdjustInventory(context.Context, int64, int, string) error { return nil }
func (s *cabinDeleteErrSvc) ListPrices(context.Context, int64) ([]domain.CabinPrice, error) {
	return nil, nil
}
func (s *cabinDeleteErrSvc) UpsertPrice(context.Context, *domain.CabinPrice) error { return nil }

type bookingAdminDeleteErrStore struct{ err error }

func (s *bookingAdminDeleteErrStore) List(context.Context, int, int) ([]domain.Booking, int64, error) {
	return nil, 0, nil
}
func (s *bookingAdminDeleteErrStore) GetByID(context.Context, int64) (*domain.Booking, error) {
	return &domain.Booking{}, nil
}
func (s *bookingAdminDeleteErrStore) UpdateStatus(context.Context, int64, string) error { return nil }
func (s *bookingAdminDeleteErrStore) Delete(context.Context, int64) error                { return s.err }

func TestDeleteConflictMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("route delete returns 409 on fk conflict", func(t *testing.T) {
		r := gin.New()
		r.DELETE("/admin/routes/:id", NewRouteHandler(&routeDeleteErrSvc{err: errors.New(fkErrMsg)}).Delete)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/admin/routes/1", nil))
		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("voyage delete returns 409 on fk conflict", func(t *testing.T) {
		r := gin.New()
		r.DELETE("/admin/voyages/:id", NewVoyageHandler(&voyageDeleteErrSvc{err: errors.New(fkErrMsg)}).Delete)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/admin/voyages/1", nil))
		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("cabin delete returns 409 on fk conflict", func(t *testing.T) {
		r := gin.New()
		r.DELETE("/admin/cabins/:id", NewCabinHandler(&cabinDeleteErrSvc{err: errors.New(fkErrMsg)}).Delete)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/admin/cabins/1", nil))
		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("booking admin delete returns 409 on fk conflict", func(t *testing.T) {
		r := gin.New()
		h := NewBookingHandler(&bookingTestSvc{}, &bookingAdminDeleteErrStore{err: errors.New(fkErrMsg)})
		r.DELETE("/admin/bookings/:id", h.AdminDelete)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/admin/bookings/1", nil))
		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
		}
	})
}
