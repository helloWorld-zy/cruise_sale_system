package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type bookingTestSvc struct {
	called bool
	err    error
}

func (s *bookingTestSvc) Create(_ context.Context, userID, voyageID, skuID int64, guests int) (*domain.Booking, error) {
	s.called = true
	if s.err != nil {
		return nil, s.err
	}
	return &domain.Booking{ID: 1, Status: "created", TotalCents: 10000}, nil
}

// TestCreateBooking 测试创建预订
func TestCreateBooking(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, "13800000000") // M-01
		c.Next()
	})

	svc := &bookingTestSvc{}
	h := NewBookingHandler(svc)
	r.POST("/api/bookings", h.Create)
	w := httptest.NewRecorder()
	body := []byte(`{"voyage_id":2,"cabin_sku_id":3,"guests":2}`)
	r.ServeHTTP(w, httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(body)))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !svc.called {
		t.Fatal("expected booking service to be called")
	}
}

// M-05: 补充缺少的边界条件测试

// TestCreateBookingMissingBody 测试缺少请求体的创建预订
func TestCreateBookingMissingBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, "1")
		c.Next()
	})
	h := NewBookingHandler(&bookingTestSvc{})
	r.POST("/api/bookings", h.Create)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("POST", "/api/bookings", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// TestCreateBookingNotAuthenticated 测试未认证时的创建预订
func TestCreateBookingNotAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewBookingHandler(&bookingTestSvc{})
	r.POST("/api/bookings", h.Create)

	w := httptest.NewRecorder()
	body := []byte(`{"voyage_id":2,"cabin_sku_id":3,"guests":2}`)
	r.ServeHTTP(w, httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(body)))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// TestCreateBookingServiceError 测试预订服务错误
func TestCreateBookingServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, "1")
		c.Next()
	})
	svc := &bookingTestSvc{err: fmt.Errorf("cabin unavailable")}
	h := NewBookingHandler(svc)
	r.POST("/api/bookings", h.Create)

	w := httptest.NewRecorder()
	body := []byte(`{"voyage_id":2,"cabin_sku_id":3,"guests":2}`)
	r.ServeHTTP(w, httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(body)))
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}
