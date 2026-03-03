package handler

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestCabinHandler_BatchStatusWritesAuditLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCabinHandler(&mockCabinSvc{})
	r.PUT("/api/v1/admin/cabins/batch-status", h.BatchUpdateStatus)

	var buf bytes.Buffer
	origin := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(origin)

	body := bytes.NewBufferString(`{"ids":[1,2],"status":0}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/cabins/batch-status", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(buf.String(), "audit bulk update cabins") {
		t.Fatalf("expected cabin audit log, got: %s", buf.String())
	}
}

func TestCabinHandler_BatchStatusRejectsOversize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCabinHandler(&mockCabinSvc{})
	r.PUT("/api/v1/admin/cabins/batch-status", h.BatchUpdateStatus)

	ids := make([]string, 0, maxBatchUpdateSize+1)
	for i := 0; i < maxBatchUpdateSize+1; i++ {
		ids = append(ids, "1")
	}
	body := bytes.NewBufferString(`{"ids":[` + strings.Join(ids, ",") + `],"status":0}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/cabins/batch-status", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
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

func TestCabinHandler_BatchSetPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCabinHandler(&mockCabinSvc{})
	r.POST("/api/v1/admin/cabins/:id/prices/batch", h.BatchSetPrice)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"start_date":"2026-03-10","end_date":"2026-03-12","occupancy":2,"price_cents":19900,"price_type":"base"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/cabins/1/prices/batch", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}

	w2 := httptest.NewRecorder()
	body2 := bytes.NewBufferString(`{"start_date":"2026-03-13","end_date":"2026-03-12","occupancy":2,"price_cents":19900,"price_type":"base"}`)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/admin/cabins/1/prices/batch", body2)
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid date range, got %d", w2.Code)
	}
}
