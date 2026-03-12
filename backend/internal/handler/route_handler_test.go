package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// fakeRouteService 实现了用于测试的 RouteService。
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

// TestRouteListHandler 测试航线列表处理器
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

// TestRouteDeleteHandler 测试航线删除处理器
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

// fakeVoyageService 实现了用于测试的 VoyageService。
type fakeVoyageService struct{}

func (f *fakeVoyageService) List(_ context.Context) ([]domain.Voyage, error) {
	return []domain.Voyage{}, nil
}
func (f *fakeVoyageService) ListPublic(_ context.Context, cruiseID int64, keyword string, page, pageSize int) ([]domain.Voyage, int64, error) {
	items := []domain.Voyage{
		{ID: 101, CruiseID: 11, Code: "VOY-101", BriefInfo: "海洋光谱号 上海-福冈-上海 5天4晚", MinPriceCents: 409900, SoldCount: 23},
		{ID: 202, CruiseID: 22, Code: "VOY-202", BriefInfo: "欧罗巴号 上海-冲绳-上海 5天4晚"},
	}
	filtered := make([]domain.Voyage, 0, len(items))
	for _, item := range items {
		if cruiseID > 0 && item.CruiseID != cruiseID {
			continue
		}
		if keyword != "" && !strings.Contains(item.BriefInfo, keyword) && !strings.Contains(item.Code, keyword) {
			continue
		}
		filtered = append(filtered, item)
	}
	if cruiseID == 0 && keyword == "" {
		return items, int64(len(items)), nil
	}
	return filtered, int64(len(filtered)), nil
}
func (f *fakeVoyageService) Create(_ context.Context, v *domain.Voyage) error { return nil }
func (f *fakeVoyageService) Update(_ context.Context, v *domain.Voyage) error { return nil }
func (f *fakeVoyageService) GetByID(_ context.Context, id int64) (*domain.Voyage, error) {
	eta := "08:30"
	etd := "17:30"
	return &domain.Voyage{
		ID:            id,
		CruiseID:      11,
		Code:          "RC101",
		BriefInfo:     "海洋光谱号 上海-福冈-上海 5天4晚",
		MinPriceCents: 409900,
		SoldCount:     23,
		FeeNote: &domain.FeeNoteContent{
			Included: []domain.ContentTextItem{{Text: "船票与指定餐食"}},
			Excluded: []domain.ContentTextItem{{Text: "签证费用", Emphasis: true}},
		},
		BookingNotice: &domain.BookingNoticeContent{
			Sections: []domain.BookingNoticeSection{{
				Key:   "documents",
				Title: "出行证件",
				Items: []domain.ContentTextItem{{Text: "请携带护照原件"}},
			}},
		},
		Itineraries: []domain.VoyageItinerary{{
			ID:                1,
			DayNo:             1,
			StopIndex:         1,
			City:              "上海",
			Summary:           "登船",
			ETATime:           &eta,
			ETDTime:           &etd,
			HasBreakfast:      true,
			HasLunch:          true,
			HasDinner:         true,
			HasAccommodation:  true,
			AccommodationText: "船上住宿",
		}},
	}, nil
}
func (f *fakeVoyageService) Delete(_ context.Context, id int64) error { return nil }

// TestVoyageListHandler 测试航次列表处理器
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

func TestVoyagePublicListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewVoyageHandler(&fakeVoyageService{})
	r.GET("/voyages", h.ListPublic)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/voyages?cruise_id=11&page=1&page_size=20", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			List  []domain.Voyage `json:"list"`
			Total int64           `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.Data.Total != 1 {
		t.Fatalf("expected total=1, got %d", body.Data.Total)
	}
	if len(body.Data.List) != 1 || body.Data.List[0].CruiseID != 11 {
		t.Fatalf("expected one filtered voyage for cruise 11, got %+v", body.Data.List)
	}
}

func TestVoyagePublicListHandler_SupportsKeyword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewVoyageHandler(&fakeVoyageService{})
	r.GET("/voyages", h.ListPublic)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/voyages?keyword=光谱&page=1&page_size=20", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			List  []domain.Voyage `json:"list"`
			Total int64           `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.Data.Total != 1 || len(body.Data.List) != 1 {
		t.Fatalf("expected one keyword matched voyage, got %+v", body.Data)
	}
	if body.Data.List[0].MinPriceCents != 409900 {
		t.Fatalf("expected backend price mapping, got %+v", body.Data.List[0])
	}
}

func TestVoyagePublicGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewVoyageHandler(&fakeVoyageService{})
	r.GET("/voyages/:id", h.Get)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/voyages/101", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data domain.Voyage `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.Data.Code != "RC101" || body.Data.MinPriceCents != 409900 {
		t.Fatalf("expected public voyage detail payload, got %+v", body.Data)
	}
}

func TestVoyagePublicGetHandler_ResolvesTemplateContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewVoyageHandler(&fakeVoyageService{})
	r.GET("/voyages/:id", h.Get)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/voyages/101", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if _, ok := body.Data["fee_note"]; !ok {
		t.Fatalf("expected fee_note in public voyage detail, got %+v", body.Data)
	}
	if _, ok := body.Data["booking_notice"]; !ok {
		t.Fatalf("expected booking_notice in public voyage detail, got %+v", body.Data)
	}
	rows, ok := body.Data["itineraries"].([]any)
	if !ok || len(rows) == 0 {
		t.Fatalf("expected itineraries in public voyage detail, got %+v", body.Data)
	}
	first, ok := rows[0].(map[string]any)
	if !ok {
		t.Fatalf("expected itinerary row map, got %+v", rows[0])
	}
	if _, ok := first["has_breakfast"]; !ok {
		t.Fatalf("expected meal flags in itinerary row, got %+v", first)
	}
	if _, ok := first["eta_time"]; !ok {
		t.Fatalf("expected eta_time in itinerary row, got %+v", first)
	}
}

func TestVoyageCreateHandler_AcceptsTemplateSelection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewVoyageHandler(&fakeVoyageService{})
	r.POST("/voyages", h.Create)
	body := `{
		"cruise_id": 11,
		"code": "RC101",
		"brief_info": "海洋光谱号 上海-福冈-上海 5天4晚",
		"depart_date": "2026-04-02T00:00:00Z",
		"return_date": "2026-04-06T00:00:00Z",
		"fee_note_template_id": 8,
		"fee_note_mode": "template",
		"booking_notice_template_id": 9,
		"booking_notice_mode": "snapshot",
		"booking_notice_content": {"sections": [{"key": "documents", "title": "出行证件", "items": [{"text": "请携带护照原件"}]}]},
		"itineraries": [
			{"day_no": 1, "stop_index": 1, "city": "上海", "summary": "登船"}
		]
	}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/voyages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var payload struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if payload.Data["fee_note_template_id"] != float64(8) {
		t.Fatalf("expected fee_note_template_id=8, got %+v", payload.Data)
	}
	if payload.Data["booking_notice_mode"] != "snapshot" {
		t.Fatalf("expected booking_notice_mode=snapshot, got %+v", payload.Data)
	}
}
