package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCabinHandler_FilteredList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCabinHandler(&mockCabinSvc{})
	r.GET("/api/v1/admin/cabins", h.FilteredList)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/admin/cabins?voyage_id=1&cabin_type_id=2&page=1&page_size=10", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCabinHandler_BatchStatusAndAlerts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCabinHandler(&mockCabinSvc{})
	r.PUT("/api/v1/admin/cabins/batch-status", h.BatchUpdateStatus)
	r.GET("/api/v1/admin/cabins/alerts", h.GetAlerts)
	r.PUT("/api/v1/admin/cabins/:id/alert-threshold", h.SetAlertThreshold)

	body := bytes.NewBufferString(`{"ids":[1,2],"status":0}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/cabins/batch-status", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/api/v1/admin/cabins/alerts", nil))
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}

	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodPut, "/api/v1/admin/cabins/1/alert-threshold", bytes.NewBufferString(`{"threshold":3}`))
	req3.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w3.Code)
	}
}

func TestCabinHandler_CategoryTree(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCabinHandler(&mockCabinSvc{})
	r.GET("/api/v1/admin/cabins/category-tree", h.CategoryTree)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/admin/cabins/category-tree", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
