package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type imageRepoForHandler struct {
	items []domain.Image
}

func (m *imageRepoForHandler) Create(ctx context.Context, img *domain.Image) error {
	_ = ctx
	img.ID = int64(len(m.items) + 1)
	m.items = append(m.items, *img)
	return nil
}

func (m *imageRepoForHandler) ListByEntity(ctx context.Context, entityType string, entityID int64) ([]domain.Image, error) {
	_ = ctx
	out := make([]domain.Image, 0)
	for _, it := range m.items {
		if it.EntityType == entityType && it.EntityID == entityID {
			out = append(out, it)
		}
	}
	return out, nil
}

func (m *imageRepoForHandler) DeleteByEntity(ctx context.Context, entityType string, entityID int64) error {
	_ = ctx
	filtered := make([]domain.Image, 0)
	for _, it := range m.items {
		if !(it.EntityType == entityType && it.EntityID == entityID) {
			filtered = append(filtered, it)
		}
	}
	m.items = filtered
	return nil
}

func (m *imageRepoForHandler) UpdateSortOrder(ctx context.Context, id int64, sortOrder int) error {
	_ = ctx
	for i := range m.items {
		if m.items[i].ID == id {
			m.items[i].SortOrder = sortOrder
			return nil
		}
	}
	return nil
}
func (m *imageRepoForHandler) ReplaceImages(ctx context.Context, entityType string, entityID int64, images []*domain.Image) error {
	_ = m.DeleteByEntity(ctx, entityType, entityID)
	for _, img := range images {
		img.ID = int64(len(m.items) + 1)
		m.items = append(m.items, *img)
	}
	return nil
}

func TestCruiseHandler_BatchUpdateStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCruiseHandler(service.NewCruiseService(&mockCruiseRepo{}, &mockCabinTypeRepo{}, &mockCompanyRepo{}))
	r.PUT("/api/v1/admin/cruises/batch-status", h.BatchUpdateStatus)

	body := bytes.NewBufferString(`{"ids":[1],"status":0}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/cruises/batch-status", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCruiseHandler_BatchUpdateStatusRejectsOversize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCruiseHandler(service.NewCruiseService(&mockCruiseRepo{}, &mockCabinTypeRepo{}, &mockCompanyRepo{}))
	r.PUT("/api/v1/admin/cruises/batch-status", h.BatchUpdateStatus)

	ids := make([]string, 0, maxBatchUpdateSize+1)
	for i := 0; i < maxBatchUpdateSize+1; i++ {
		ids = append(ids, "1")
	}
	body := bytes.NewBufferString(`{"ids":[` + strings.Join(ids, ",") + `],"status":0}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/cruises/batch-status", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCruiseHandler_BatchUpdateStatusWritesAuditLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCruiseHandler(service.NewCruiseService(&mockCruiseRepo{}, &mockCabinTypeRepo{}, &mockCompanyRepo{}))
	r.PUT("/api/v1/admin/cruises/batch-status", h.BatchUpdateStatus)

	var buf bytes.Buffer
	origin := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(origin)

	body := bytes.NewBufferString(`{"ids":[1],"status":0}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/cruises/batch-status", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(buf.String(), "audit bulk update cruises") {
		t.Fatalf("expected cruise audit log, got: %s", buf.String())
	}
}

func TestImageHandler_SaveAndList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := &imageRepoForHandler{}
	h := NewImageHandler(service.NewImageService(repo))
	r.POST("/api/v1/admin/images", h.Save)
	r.GET("/api/v1/admin/images", h.List)

	postBody := map[string]any{
		"entity_type": "cruise",
		"entity_id":   1,
		"images": []map[string]any{
			{"url": "https://img/a.jpg", "sort_order": 1, "is_primary": true},
		},
	}
	data, _ := json.Marshal(postBody)
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/images", bytes.NewReader(data))
	postReq.Header.Set("Content-Type", "application/json")
	postResp := httptest.NewRecorder()
	r.ServeHTTP(postResp, postReq)
	if postResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", postResp.Code, postResp.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/images?entity_type=cruise&entity_id=1", nil)
	listResp := httptest.NewRecorder()
	r.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", listResp.Code, listResp.Body.String())
	}
}

func TestFacilityHandler_GetAndUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewFacilityHandler(service.NewFacilityService(&mockFacilityRepo{}))
	r.GET("/api/v1/admin/facilities/:id", h.Get)
	r.PUT("/api/v1/admin/facilities/:id", h.Update)

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/facilities/1", nil)
	getResp := httptest.NewRecorder()
	r.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", getResp.Code, getResp.Body.String())
	}

	putBody := bytes.NewBufferString(`{"cruise_id":1,"category_id":1,"name":"spa"}`)
	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/admin/facilities/1", putBody)
	putReq.Header.Set("Content-Type", "application/json")
	putResp := httptest.NewRecorder()
	r.ServeHTTP(putResp, putReq)
	if putResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", putResp.Code, putResp.Body.String())
	}
}

func TestFacilityCategoryHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewFacilityCategoryHandler(service.NewFacilityCategoryService(&mockFacilityCategoryRepo{}))
	r.PUT("/api/v1/admin/facility-categories/:id", h.Update)

	body := bytes.NewBufferString(`{"name":"餐饮","icon":"fork","status":1,"sort_order":1}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/facility-categories/1", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", resp.Code, resp.Body.String())
	}
}

type captureCruiseRepo struct {
	updated *domain.Cruise
}

func (m *captureCruiseRepo) Create(ctx context.Context, cruise *domain.Cruise) error { return nil }
func (m *captureCruiseRepo) Update(ctx context.Context, cruise *domain.Cruise) error {
	m.updated = &domain.Cruise{ID: cruise.ID, CompanyID: cruise.CompanyID, Name: cruise.Name}
	return nil
}
func (m *captureCruiseRepo) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	if id != 1 {
		return nil, errors.New("not found")
	}
	return &domain.Cruise{ID: 1, CompanyID: 1, Name: "old"}, nil
}
func (m *captureCruiseRepo) List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	return nil, 0, nil
}
func (m *captureCruiseRepo) ListPublic(ctx context.Context, companyID int64, keyword string, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	return nil, 0, nil
}
func (m *captureCruiseRepo) Delete(ctx context.Context, id int64) error { return nil }

type passCabinTypeRepo struct{}

func (m *passCabinTypeRepo) Create(ctx context.Context, c *domain.CabinType) error { return nil }
func (m *passCabinTypeRepo) Update(ctx context.Context, c *domain.CabinType) error { return nil }
func (m *passCabinTypeRepo) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	return &domain.CabinType{ID: id}, nil
}
func (m *passCabinTypeRepo) ListByCruise(ctx context.Context, cruiseID int64, page, pageSize int) ([]domain.CabinType, int64, error) {
	return nil, 0, nil
}
func (m *passCabinTypeRepo) Delete(ctx context.Context, id int64) error { return nil }

type passCompanyRepo struct{}

func (m *passCompanyRepo) Create(ctx context.Context, company *domain.CruiseCompany) error {
	return nil
}
func (m *passCompanyRepo) Update(ctx context.Context, company *domain.CruiseCompany) error {
	return nil
}
func (m *passCompanyRepo) GetByID(ctx context.Context, id int64) (*domain.CruiseCompany, error) {
	return &domain.CruiseCompany{ID: id}, nil
}
func (m *passCompanyRepo) List(ctx context.Context, keyword string, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	return nil, 0, nil
}
func (m *passCompanyRepo) ListPublic(ctx context.Context, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	return nil, 0, nil
}
func (m *passCompanyRepo) Delete(ctx context.Context, id int64) error { return nil }

func TestCruiseHandler_Update_ChangesCompanyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cruiseRepo := &captureCruiseRepo{}
	svc := service.NewCruiseService(cruiseRepo, &passCabinTypeRepo{}, &passCompanyRepo{})
	h := NewCruiseHandler(svc)
	r.PUT("/api/v1/admin/cruises/:id", h.Update)

	body := bytes.NewBufferString(`{"company_id":2,"name":"new-name"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/cruises/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if cruiseRepo.updated == nil {
		t.Fatal("expected cruise update to be called")
	}
	if cruiseRepo.updated.CompanyID != 2 {
		t.Fatalf("expected updated company_id=2, got %d", cruiseRepo.updated.CompanyID)
	}
}
