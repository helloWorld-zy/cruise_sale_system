package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/repository"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// --- mock BookingAdminStore ---
type mockBookingAdminStore struct{}

func (m *mockBookingAdminStore) List(_ context.Context, page, pageSize int) ([]domain.Booking, int64, error) {
	if page == 99 {
		return nil, 0, errors.New("list error")
	}
	return []domain.Booking{{ID: 1}}, 1, nil
}
func (m *mockBookingAdminStore) ListWithFilter(_ context.Context, filter repository.BookingFilter, page, pageSize int) ([]domain.Booking, int64, error) {
	return []domain.Booking{{ID: 1}}, 1, nil
}
func (m *mockBookingAdminStore) GetByID(_ context.Context, id int64) (*domain.Booking, error) {
	if id == 99 {
		return nil, errors.New("not found")
	}
	return &domain.Booking{ID: id}, nil
}
func (m *mockBookingAdminStore) TransitionStatus(_ context.Context, id int64, status string, operatorID int64, remark string) error {
	_ = operatorID
	_ = remark
	if id == 99 {
		return errors.New("update error")
	}
	return nil
}
func (m *mockBookingAdminStore) Delete(_ context.Context, id int64) error {
	if id == 99 {
		return errors.New("delete error")
	}
	return nil
}

func doJSONReq(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf []byte
	if body != nil {
		buf, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ========== Booking AdminList / AdminGet / AdminUpdate / AdminDelete ==========

func TestBookingHandler_AdminList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewBookingHandler(&mockBookingSvc{}, &mockBookingAdminStore{})
	r.GET("/bookings", h.AdminList)

	w := doJSONReq(r, "GET", "/bookings?page=1&page_size=10", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// adminStore == nil
	h2 := NewBookingHandler(&mockBookingSvc{})
	r2 := gin.New()
	r2.GET("/bookings", h2.AdminList)
	w2 := doJSONReq(r2, "GET", "/bookings", nil)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestBookingHandler_AdminGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewBookingHandler(&mockBookingSvc{}, &mockBookingAdminStore{})
	r.GET("/bookings/:id", h.AdminGet)

	w := doJSONReq(r, "GET", "/bookings/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "GET", "/bookings/99", nil)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	w3 := doJSONReq(r, "GET", "/bookings/x", nil)
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	w4 := doJSONReq(r, "GET", "/bookings/0", nil)
	assert.Equal(t, http.StatusBadRequest, w4.Code)

	// adminStore nil
	h2 := NewBookingHandler(&mockBookingSvc{})
	r2 := gin.New()
	r2.GET("/bookings/:id", h2.AdminGet)
	w5 := doJSONReq(r2, "GET", "/bookings/1", nil)
	assert.Equal(t, http.StatusInternalServerError, w5.Code)
}

func TestBookingHandler_AdminUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewBookingHandler(&mockBookingSvc{}, &mockBookingAdminStore{})
	r.PUT("/bookings/:id", h.AdminUpdate)

	w := doJSONReq(r, "PUT", "/bookings/1", map[string]string{"status": "paid"})
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "PUT", "/bookings/99", map[string]string{"status": "paid"})
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

	w3 := doJSONReq(r, "PUT", "/bookings/x", nil)
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	w4 := doJSONReq(r, "PUT", "/bookings/1", map[string]int{"status": 1})
	assert.Equal(t, http.StatusBadRequest, w4.Code)

	// adminStore nil
	h2 := NewBookingHandler(&mockBookingSvc{})
	r2 := gin.New()
	r2.PUT("/bookings/:id", h2.AdminUpdate)
	w5 := doJSONReq(r2, "PUT", "/bookings/1", map[string]string{"status": "paid"})
	assert.Equal(t, http.StatusInternalServerError, w5.Code)
}

func TestBookingHandler_AdminDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewBookingHandler(&mockBookingSvc{}, &mockBookingAdminStore{})
	r.DELETE("/bookings/:id", h.AdminDelete)

	w := doJSONReq(r, "DELETE", "/bookings/1", nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	w2 := doJSONReq(r, "DELETE", "/bookings/99", nil)
	assert.True(t, w2.Code >= 400)

	w3 := doJSONReq(r, "DELETE", "/bookings/x", nil)
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	// adminStore nil
	h2 := NewBookingHandler(&mockBookingSvc{})
	r2 := gin.New()
	r2.DELETE("/bookings/:id", h2.AdminDelete)
	w4 := doJSONReq(r2, "DELETE", "/bookings/1", nil)
	assert.Equal(t, http.StatusInternalServerError, w4.Code)
}

// ========== CabinHandler.List / CabinHandler.Get ==========

func TestCabinHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCabinHandler(&mockCabinSvc{})
	r.GET("/cabins/list", h.List)

	w := doJSONReq(r, "GET", "/cabins/list?voyage_id=1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "GET", "/cabins/list?voyage_id=99", nil)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestCabinHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCabinHandler(&mockCabinSvc{})
	r.GET("/cabins/:id", h.Get)

	w := doJSONReq(r, "GET", "/cabins/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "GET", "/cabins/99", nil)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	w3 := doJSONReq(r, "GET", "/cabins/x", nil)
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	w4 := doJSONReq(r, "GET", "/cabins/0", nil)
	assert.Equal(t, http.StatusBadRequest, w4.Code)
}

// ========== CruiseHandler.Get ==========

func TestCruiseHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewCruiseService(&mockCruiseRepo{}, &mockCabinTypeRepo{}, &mockCompanyRepo{})
	h := NewCruiseHandler(svc)
	r.GET("/cruises/:id", h.Get)

	w := doJSONReq(r, "GET", "/cruises/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "GET", "/cruises/99", nil)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	w3 := doJSONReq(r, "GET", "/cruises/x", nil)
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	w4 := doJSONReq(r, "GET", "/cruises/0", nil)
	assert.Equal(t, http.StatusBadRequest, w4.Code)
}

// ========== CruiseHandler.BatchUpdateStatus ==========

func TestCruiseHandler_BatchUpdateStatus_ExtraCases(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewCruiseService(&mockCruiseRepo{}, &mockCabinTypeRepo{}, &mockCompanyRepo{})
	h := NewCruiseHandler(svc)
	r.PUT("/cruises/batch-status", h.BatchUpdateStatus)

	// empty ids
	w := doJSONReq(r, "PUT", "/cruises/batch-status", map[string]interface{}{"ids": []int64{}, "status": 1})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// valid
	w2 := doJSONReq(r, "PUT", "/cruises/batch-status", map[string]interface{}{"ids": []int64{1}, "status": 1})
	assert.Equal(t, http.StatusOK, w2.Code)

	// bad body
	w3 := doJSONReq(r, "PUT", "/cruises/batch-status", "bad")
	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

// ========== RouteHandler.Get ==========

func TestRouteHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewRouteHandler(&mockRouteSvc{})
	r.GET("/routes/:id", h.Get)

	w := doJSONReq(r, "GET", "/routes/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "GET", "/routes/99", nil)
	assert.True(t, w2.Code >= 400) // not found

	w3 := doJSONReq(r, "GET", "/routes/x", nil)
	assert.True(t, w3.Code >= 400) // bad id
}

// ========== VoyageHandler.Get ==========

func TestVoyageHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewVoyageHandler(&mockVoyageSvc{})
	r.GET("/voyages/:id", h.Get)

	w := doJSONReq(r, "GET", "/voyages/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "GET", "/voyages/99", nil)
	assert.True(t, w2.Code >= 400)

	w3 := doJSONReq(r, "GET", "/voyages/x", nil)
	assert.True(t, w3.Code >= 400)
}

// ========== FacilityCategoryHandler.Update (distinct from cruise_handler_extended_test) ==========

func TestFacilityCategoryHandler_UpdateP06(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewFacilityCategoryService(&mockFacilityCategoryRepo{})
	h := NewFacilityCategoryHandler(svc)
	r.PUT("/fc/:id", h.Update)

	w := doJSONReq(r, "PUT", "/fc/1", map[string]interface{}{"name": "ok"})
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "PUT", "/fc/x", nil)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	w3 := doJSONReq(r, "PUT", "/fc/99", map[string]interface{}{"name": "ok"})
	assert.True(t, w3.Code >= 400)

	w4 := doJSONReq(r, "PUT", "/fc/1", map[string]interface{}{"name": "error"})
	assert.True(t, w4.Code >= 400)

	w5 := doJSONReq(r, "PUT", "/fc/1", map[string]int{"name": 1})
	assert.Equal(t, http.StatusBadRequest, w5.Code)
}

// ========== FacilityHandler.Get / .Update ==========

func TestFacilityHandler_GetP06(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewFacilityService(&mockFacilityRepo{})
	h := NewFacilityHandler(svc)
	r.GET("/facilities/:id", h.Get)

	w := doJSONReq(r, "GET", "/facilities/1", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "GET", "/facilities/99", nil)
	assert.True(t, w2.Code >= 400)

	w3 := doJSONReq(r, "GET", "/facilities/x", nil)
	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

func TestFacilityHandler_UpdateP06(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := service.NewFacilityService(&mockFacilityRepo{})
	h := NewFacilityHandler(svc)
	r.PUT("/facilities/:id", h.Update)

	w := doJSONReq(r, "PUT", "/facilities/1", map[string]interface{}{"cruise_id": 1, "category_id": 1, "name": "ok"})
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := doJSONReq(r, "PUT", "/facilities/x", nil)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	w3 := doJSONReq(r, "PUT", "/facilities/1", map[string]int{"name": 1})
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	w4 := doJSONReq(r, "PUT", "/facilities/1", map[string]interface{}{"cruise_id": 1, "category_id": 1, "name": "error"})
	assert.True(t, w4.Code >= 400)
}
