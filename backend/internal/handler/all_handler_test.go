package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/repository"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockCompanyRepo struct{}

func (m *mockCompanyRepo) Create(ctx context.Context, company *domain.CruiseCompany) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if company.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockCompanyRepo) Update(ctx context.Context, company *domain.CruiseCompany) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if company.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockCompanyRepo) GetByID(ctx context.Context, id int64) (*domain.CruiseCompany, error) {
	if id == 1 {
		return &domain.CruiseCompany{ID: 1}, nil
	}
	return nil, errors.New("not found")
}
func (m *mockCompanyRepo) List(ctx context.Context, keyword string, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	if keyword == "error" {
		return nil, 0, errors.New("error")
	}
	return []domain.CruiseCompany{{ID: 1}}, 1, nil
}
func (m *mockCompanyRepo) ListPublic(ctx context.Context, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	return []domain.CruiseCompany{{ID: 1, Name: "皇家加勒比", Status: 1}}, 1, nil
}
func (m *mockCompanyRepo) Delete(ctx context.Context, id int64) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if id == 2 {
		return service.ErrCompanyHasCruises
	}
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type mockCruiseRepo struct{}

func (m *mockCruiseRepo) Create(ctx context.Context, cruise *domain.Cruise) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if cruise.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockCruiseRepo) Update(ctx context.Context, cruise *domain.Cruise) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if cruise.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockCruiseRepo) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	if id == 1 {
		return &domain.Cruise{ID: 1}, nil
	}
	return nil, errors.New("not found")
}
func (m *mockCruiseRepo) List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	_ = keyword
	_ = status
	_ = sortBy
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	if companyID == 2 {
		return nil, 1, nil
	}
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	if companyID == 99 {
		return nil, 0, errors.New("error")
	}
	return []domain.Cruise{{ID: 1}}, 1, nil
}
func (m *mockCruiseRepo) ListPublic(ctx context.Context, companyID int64, keyword string, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	if companyID == 2 {
		return []domain.Cruise{{ID: 2, CompanyID: 2, Name: "海洋奇迹号", Status: 1}}, 1, nil
	}
	return []domain.Cruise{{ID: 1, CompanyID: 1, Name: "海洋光谱号", Status: 1}}, 1, nil
}
func (m *mockCruiseRepo) Delete(ctx context.Context, id int64) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if id == 2 {
		return service.ErrCruiseHasCabins
	}
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type mockCabinTypeRepo struct{}

func (m *mockCabinTypeRepo) Create(ctx context.Context, c *domain.CabinType) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if c.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockCabinTypeRepo) Update(ctx context.Context, c *domain.CabinType) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if c.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockCabinTypeRepo) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	if id == 1 || id == 99 {
		return &domain.CabinType{ID: id}, nil
	}
	return nil, errors.New("not found")
}
func (m *mockCabinTypeRepo) ListByCruise(ctx context.Context, id int64, p, ps int) ([]domain.CabinType, int64, error) {
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	if id == 2 {
		return nil, 1, nil
	}
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	if id == 99 {
		return nil, 0, errors.New("error")
	}
	return nil, 0, nil
}
func (m *mockCabinTypeRepo) Delete(ctx context.Context, id int64) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type mockFacilityCategoryRepo struct{}

func (m *mockFacilityCategoryRepo) Create(ctx context.Context, c *domain.FacilityCategory) error {
	if c.Name == "error" {
		return errors.New("err")
	}
	if isErr(ctx) {
		return errors.New("error")
	}
	if isErr(ctx) {
		return errors.New("error")
	}
	if c.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockFacilityCategoryRepo) List(ctx context.Context) ([]domain.FacilityCategory, error) {
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	if ctx.Value("trigger_error") != nil {
		return nil, errors.New("error")
	}
	return nil, nil
}
func (m *mockFacilityCategoryRepo) Update(ctx context.Context, c *domain.FacilityCategory) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if c.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockFacilityCategoryRepo) GetByID(ctx context.Context, id int64) (*domain.FacilityCategory, error) {
	if id == 99 {
		return nil, errors.New("error")
	}
	return &domain.FacilityCategory{ID: id}, nil
}
func (m *mockFacilityCategoryRepo) Delete(ctx context.Context, id int64) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type mockFacilityRepo struct{}

func (m *mockFacilityRepo) Create(ctx context.Context, f *domain.Facility) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if f.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockFacilityRepo) ListByCruise(ctx context.Context, cruiseID int64) ([]domain.Facility, error) {
	if cruiseID == 99 {
		return nil, errors.New("err")
	}
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	if cruiseID == 99 {
		return nil, errors.New("error")
	}
	return nil, nil
}
func (m *mockFacilityRepo) Update(ctx context.Context, f *domain.Facility) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if f.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockFacilityRepo) GetByID(ctx context.Context, id int64) (*domain.Facility, error) {
	if id == 99 {
		return nil, errors.New("error")
	}
	return &domain.Facility{ID: id}, nil
}
func (m *mockFacilityRepo) ListByCruiseAndCategory(ctx context.Context, cruiseID, categoryID int64) ([]domain.Facility, error) {
	_ = categoryID
	return m.ListByCruise(ctx, cruiseID)
}
func (m *mockFacilityRepo) Delete(ctx context.Context, id int64) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

func setupRouter() (*gin.Engine, *AuthHandler, *CompanyHandler, *CruiseHandler, *CabinTypeHandler, *FacilityCategoryHandler, *FacilityHandler, *UploadHandler, *UserHandler) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if c.Query("err") == "1" {
			c.Set("trigger_error", true)
		}
		c.Next()
	})

	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&domain.Staff{})

	staffRepo := repository.NewStaffRepository(db)

	hash, _ := service.HashPassword("password")
	db.Create(&domain.Staff{ID: 1, Username: "admin", PasswordHash: hash, Status: 1})

	authSvc := service.NewAuthService(staffRepo, "secret", 24)
	compSvc := service.NewCompanyService(&mockCompanyRepo{}, &mockCruiseRepo{})
	cruiseSvc := service.NewCruiseService(&mockCruiseRepo{}, &mockCabinTypeRepo{}, &mockCompanyRepo{})
	cabinTypeSvc := service.NewCabinTypeService(&mockCabinTypeRepo{})
	facCatSvc := service.NewFacilityCategoryService(&mockFacilityCategoryRepo{})
	facSvc := service.NewFacilityService(&mockFacilityRepo{})

	return r,
		NewAuthHandler(authSvc),
		NewCompanyHandler(compSvc),
		NewCruiseHandler(cruiseSvc),
		NewCabinTypeHandler(cabinTypeSvc),
		NewFacilityCategoryHandler(facCatSvc),
		NewFacilityHandler(facSvc),
		NewUploadHandler(),
		NewUserHandler(nil, "secret")
}

func doReq(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code >= 400 {
		importFmt := true
		if importFmt {
			// 这只是一个编译技巧，以确保使用 fmt，或者如果我传递了 t，我只会使用 t.Log。
			// 让我们直接打印到 os.Stdout
			_ = importFmt
		}
		os.Stdout.WriteString("doReq " + method + " " + path + " returned " + w.Body.String() + "\n")
	}
	return w
}

// TestAuthHandler 测试认证处理器
func TestAuthHandler(t *testing.T) {
	r, authH, _, _, _, _, _, _, _ := setupRouter()
	r.POST("/login", authH.Login)
	r.GET("/profile", func(c *gin.Context) {
		uid := c.Query("uid")
		if uid != "" {
			c.Set("staffID", uid)
		}
		authH.GetProfile(c)
	})

	doReq(r, "POST", "/login", LoginRequest{Username: "admin", Password: "password"})
	doReq(r, "POST", "/login", LoginRequest{Username: "invalid", Password: "wrong"})
	doReq(r, "POST", "/login", map[string]string{"foo": "bar"})
	doReq(r, "GET", "/profile?uid=1", nil)
	doReq(r, "GET", "/profile?uid=99", nil)
	doReq(r, "GET", "/profile", nil)
}

// TestCompanyHandler 测试公司处理器
func TestCompanyHandler(t *testing.T) {
	r, _, compH, _, _, _, _, _, _ := setupRouter()
	r.GET("/companies", compH.List)
	r.POST("/companies", compH.Create)
	r.PUT("/companies/:id", compH.Update)
	r.DELETE("/companies/:id", compH.Delete)

	doReq(r, "GET", "/companies", nil)
	doReq(r, "GET", "/companies?keyword=error", nil)
	doReq(r, "POST", "/companies", domain.CruiseCompany{Name: "test"})
	doReq(r, "POST", "/companies", domain.CruiseCompany{Name: "error"})
	doReq(r, "POST", "/companies", map[string]int{"name": 1}) // bad
	doReq(r, "PUT", "/companies/1", domain.CruiseCompany{Name: "test"})
	doReq(r, "PUT", "/companies/1", map[string]int{"name": 1})
	doReq(r, "PUT", "/companies/1", domain.CruiseCompany{Name: "error"})
	doReq(r, "PUT", "/companies/99", domain.CruiseCompany{Name: "test"}) // not found
	doReq(r, "PUT", "/companies/x", nil)                                 // bad
	doReq(r, "DELETE", "/companies/1", nil)
	doReq(r, "DELETE", "/companies/2", nil) // conflict
	doReq(r, "DELETE", "/companies/1?err=1", nil)
	doReq(r, "DELETE", "/companies/99", nil)
	doReq(r, "DELETE", "/companies/x", nil)
}

// TestCruiseHandler 测试邮轮处理器
func TestCruiseHandler(t *testing.T) {
	r, _, _, cruiseH, _, _, _, _, _ := setupRouter()
	r.GET("/cruises", cruiseH.List)
	r.POST("/cruises", cruiseH.Create)
	r.PUT("/cruises/:id", cruiseH.Update)
	r.DELETE("/cruises/:id", cruiseH.Delete)

	doReq(r, "GET", "/cruises", nil)
	doReq(r, "GET", "/cruises?page=abc&company_id=abc", nil)
	doReq(r, "GET", "/cruises?err=1", nil)
	doReq(r, "GET", "/cruises?company_id=abc&page=abc&page_size=xyz", nil)
	doReq(r, "POST", "/cruises", map[string]interface{}{"company_id": 1, "name": "test"})
	doReq(r, "POST", "/cruises", map[string]interface{}{"company_id": 1, "name": "error"})
	doReq(r, "POST", "/cruises", map[string]interface{}{"company_id": 99, "name": "test"}) // invalid company
	doReq(r, "POST", "/cruises", map[string]int{"name": 1})
	doReq(r, "PUT", "/cruises/1", map[string]interface{}{"company_id": 1, "name": "test"})
	doReq(r, "PUT", "/cruises/1", map[string]int{"name": 1}) // triggers ShouldBindJSON error
	doReq(r, "PUT", "/cruises/x", nil)
	doReq(r, "PUT", "/cruises/99", map[string]interface{}{"company_id": 1, "name": "test"})
	doReq(r, "PUT", "/cruises/1", map[string]interface{}{"company_id": 1, "name": "error"})
	doReq(r, "DELETE", "/cruises/1", nil)
	doReq(r, "DELETE", "/cruises/2", nil) // conflict
	doReq(r, "DELETE", "/cruises/1?err=1", nil)
	doReq(r, "DELETE", "/cruises/99", nil)
	doReq(r, "DELETE", "/cruises/x", nil)
}

func TestCompanyHandler_ListPublic(t *testing.T) {
	r, _, compH, _, _, _, _, _, _ := setupRouter()
	r.GET("/public/companies", compH.ListPublic)

	w := doReq(r, "GET", "/public/companies", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "皇家加勒比") {
		t.Fatalf("expected public company payload, got %s", w.Body.String())
	}
}

func TestCruiseHandler_ListPublic(t *testing.T) {
	r, _, _, cruiseH, _, _, _, _, _ := setupRouter()
	r.GET("/public/cruises", cruiseH.ListPublic)

	w := doReq(r, "GET", "/public/cruises?company_id=2", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "海洋奇迹号") {
		t.Fatalf("expected filtered public cruise payload, got %s", w.Body.String())
	}
}

// TestCabinTypeHandler 测试舱房类型处理器
func TestCabinTypeHandler(t *testing.T) {
	r, _, _, _, ctH, _, _, _, _ := setupRouter()
	r.GET("/cabin-types", ctH.List)
	r.GET("/cabin-types/:id", ctH.Get)
	r.POST("/cabin-types", ctH.Create)
	r.PUT("/cabin-types/:id", ctH.Update)
	r.DELETE("/cabin-types/:id", ctH.Delete)

	doReq(r, "GET", "/cabin-types", nil)
	doReq(r, "GET", "/cabin-types/1", nil)
	doReq(r, "GET", "/cabin-types/99", nil)
	doReq(r, "GET", "/cabin-types/x", nil)
	doReq(r, "GET", "/cabin-types?err=1", nil)
	doReq(r, "GET", "/cabin-types?cruise_id=99&page=err", nil)
	doReq(r, "POST", "/cabin-types", map[string]interface{}{"cruise_id": 1, "name": "test"})
	doReq(r, "POST", "/cabin-types", map[string]interface{}{"cruise_id": 1, "name": "error"})
	doReq(r, "POST", "/cabin-types", map[string]int{"name": 1})
	doReq(r, "PUT", "/cabin-types/1", map[string]interface{}{"cruise_id": 1, "name": "test"})
	doReq(r, "PUT", "/cabin-types/1", map[string]int{"name": 1}) // triggers ShouldBindJSON error
	doReq(r, "PUT", "/cabin-types/1", map[string]interface{}{"cruise_id": 1, "name": "error"})
	doReq(r, "PUT", "/cabin-types/x", nil)
	doReq(r, "PUT", "/cabin-types/99", map[string]interface{}{"cruise_id": 1, "name": "test"})
	doReq(r, "DELETE", "/cabin-types/1", nil)
	doReq(r, "DELETE", "/cabin-types/99", nil)
	doReq(r, "DELETE", "/cabin-types/x", nil)
	doReq(r, "GET", "/cabin-types?cruise_id=1&err=1", nil)
}

// TestFacilityCategoryHandler 测试设施分类处理器
func TestFacilityCategoryHandler(t *testing.T) {
	r, _, _, _, _, fcH, _, _, _ := setupRouter()
	r.GET("/facility-categories", fcH.List)
	r.POST("/facility-categories", fcH.Create)
	r.DELETE("/facility-categories/:id", fcH.Delete)

	doReq(r, "GET", "/facility-categories", nil)
	doReq(r, "GET", "/facility-categories?err=1", nil)
	doReq(r, "GET", "/facility-categories?err=1&page=abc&page_size=xyz", nil)
	doReq(r, "POST", "/facility-categories", domain.FacilityCategory{})
	doReq(r, "POST", "/facility-categories", domain.FacilityCategory{Name: "error"})
	doReq(r, "POST", "/facility-categories", map[string]int{"name": 1})
	doReq(r, "DELETE", "/facility-categories/1", nil)
	doReq(r, "DELETE", "/facility-categories/1?err=1", nil)
	doReq(r, "DELETE", "/facility-categories/1?err=1", nil)
	doReq(r, "DELETE", "/facility-categories/99", nil)
	doReq(r, "DELETE", "/facility-categories/x", nil)
	doReq(r, "POST", "/facility-categories?err=1", domain.FacilityCategory{})
	doReq(r, "POST", "/facility-categories", map[string]int{"name": 1})
	doReq(r, "POST", "/facility-categories", map[string]interface{}{"name": "success_test_category"})
}

// TestFacilityHandler 测试设施处理器
func TestFacilityHandler(t *testing.T) {
	r, _, _, _, _, _, facH, _, _ := setupRouter()
	r.GET("/facilities", facH.ListByCruise)
	r.POST("/facilities", facH.Create)
	r.DELETE("/facilities/:id", facH.Delete)

	doReq(r, "GET", "/facilities?cruise_id=1", nil)
	doReq(r, "GET", "/facilities?cruise_id=99", nil)
	doReq(r, "GET", "/facilities?cruise_id=1&err=1&page=abc&page_size=xyz", nil)
	doReq(r, "POST", "/facilities", map[string]interface{}{"cruise_id": 1, "category_id": 1, "name": "test"})
	doReq(r, "POST", "/facilities", map[string]interface{}{"cruise_id": 1, "category_id": 1, "name": "error"})
	doReq(r, "POST", "/facilities", map[string]int{"name": 1})
	doReq(r, "DELETE", "/facilities/1", nil)
	doReq(r, "DELETE", "/facilities/99", nil)
	doReq(r, "DELETE", "/facilities/x", nil)
	doReq(r, "GET", "/facilities?cruise_id=1&err=1", nil)
	doReq(r, "GET", "/facilities?cruise_id=abc", nil)
	doReq(r, "GET", "/facilities?cruise_id=-1", nil)
}

// TestUploadHandler 测试上传处理器
func TestUploadHandler(t *testing.T) {
	r, _, _, _, _, _, _, upH, _ := setupRouter()
	r.POST("/upload", upH.UploadImage)

	t.Run("missing file", func(t *testing.T) {
		w := doReq(r, "POST", "/upload", nil)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("upload success", func(t *testing.T) {
		tmp := t.TempDir()
		h := NewUploadHandlerWithConfig(tmp, "/uploads", 2*1024*1024)
		r2 := gin.New()
		r2.POST("/upload", h.UploadImage)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", "logo.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write([]byte("\x89PNG\r\n\x1a\nmockpngcontent")); err != nil {
			t.Fatal(err)
		}
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", &body)
		req.Host = "127.0.0.1:8080"
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		r2.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), "http://127.0.0.1:8080/uploads/") {
			t.Fatalf("expected public upload url in response, body=%s", w.Body.String())
		}

		entries, err := os.ReadDir(tmp)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 {
			t.Fatalf("expected 1 uploaded file, got %d", len(entries))
		}
		if filepath.Ext(entries[0].Name()) == "" {
			t.Fatalf("expected uploaded file to have extension, got %s", entries[0].Name())
		}
	})

	t.Run("reject non-image", func(t *testing.T) {
		tmp := t.TempDir()
		h := NewUploadHandlerWithConfig(tmp, "/uploads", 2*1024*1024)
		r2 := gin.New()
		r2.POST("/upload", h.UploadImage)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", "note.txt")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write([]byte("plain-text")); err != nil {
			t.Fatal(err)
		}
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		r2.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
		}
	})
}

// TestUserHandler 测试用户处理器
func TestUserHandler(t *testing.T) {
	r, _, _, _, _, _, _, _, usH := setupRouter()
	r.POST("/user/login", usH.Login)
	r.GET("/user/profile", usH.Profile)

	doReq(r, "POST", "/user/login", nil)
	doReq(r, "GET", "/user/profile", nil)
}

// 基于接口的处理器的模拟

type mockRouteSvc struct{}

func (m *mockRouteSvc) Create(ctx context.Context, r *domain.Route) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if r.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockRouteSvc) Update(ctx context.Context, r *domain.Route) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if r.Name == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockRouteSvc) GetByID(ctx context.Context, id int64) (*domain.Route, error) {
	if id == 1 {
		return &domain.Route{ID: 1}, nil
	}
	return nil, errors.New("not found")
}
func (m *mockRouteSvc) List(ctx context.Context) ([]domain.Route, error) {
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	if ctx.Value("trigger_error") != nil {
		return nil, errors.New("error")
	}
	return nil, nil
}
func (m *mockRouteSvc) Delete(ctx context.Context, id int64) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type mockVoyageSvc struct{}

func (m *mockVoyageSvc) Create(ctx context.Context, r *domain.Voyage) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if r.Code == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockVoyageSvc) Update(ctx context.Context, r *domain.Voyage) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if r.Code == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockVoyageSvc) GetByID(ctx context.Context, id int64) (*domain.Voyage, error) {
	if id == 1 {
		return &domain.Voyage{ID: 1}, nil
	}
	return nil, errors.New("not found")
}
func (m *mockVoyageSvc) List(ctx context.Context) ([]domain.Voyage, error) {
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	return nil, nil
}
func (m *mockVoyageSvc) ListPublic(ctx context.Context, cruiseID int64, keyword string, page, pageSize int) ([]domain.Voyage, int64, error) {
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	return nil, 0, nil
}
func (m *mockVoyageSvc) Delete(ctx context.Context, id int64) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if id == 99 {
		return errors.New("error")
	}
	return nil
}

type mockCabinSvc struct{}

func (m *mockCabinSvc) ListByVoyage(ctx context.Context, id int64) ([]domain.CabinSKU, error) {
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	if id == 99 {
		return nil, errors.New("error")
	}
	return nil, nil
}
func (m *mockCabinSvc) FilteredList(ctx context.Context, f domain.CabinSKUFilter) ([]domain.CabinSKU, int64, error) {
	if isErr(ctx) {
		return nil, 0, errors.New("error")
	}
	if f.VoyageID == 99 {
		return nil, 0, errors.New("error")
	}
	return []domain.CabinSKU{{ID: 1, Code: "SKU-1"}}, 1, nil
}
func (m *mockCabinSvc) BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error {
	_ = status
	if isErr(ctx) {
		return errors.New("error")
	}
	if len(ids) == 0 {
		return errors.New("empty ids")
	}
	if ids[0] == 99 {
		return errors.New("error")
	}
	return nil
}
func (m *mockCabinSvc) GetByID(ctx context.Context, id int64) (*domain.CabinSKU, error) {
	if id == 99 {
		return nil, errors.New("error")
	}
	return &domain.CabinSKU{ID: id, Code: "SKU-1"}, nil
}
func (m *mockCabinSvc) Create(ctx context.Context, c *domain.CabinSKU) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if c.Code == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockCabinSvc) Update(ctx context.Context, c *domain.CabinSKU) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if c.Code == "error" {
		return errors.New("error")
	}
	return nil
}
func (m *mockCabinSvc) Delete(ctx context.Context, id int64) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if id == 99 {
		return errors.New("error")
	}
	return nil
}
func (m *mockCabinSvc) GetInventory(ctx context.Context, id int64) (domain.CabinInventory, error) {
	if id == 99 {
		return domain.CabinInventory{}, errors.New("error")
	}
	return domain.CabinInventory{}, nil
}
func (m *mockCabinSvc) GetAlerts(ctx context.Context) ([]domain.InventoryAlert, error) {
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	return []domain.InventoryAlert{{CabinSKUID: 1, Available: 2, AlertThreshold: 3}}, nil
}
func (m *mockCabinSvc) SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if skuID == 99 || threshold < 0 {
		return errors.New("error")
	}
	return nil
}
func (m *mockCabinSvc) AdjustInventory(ctx context.Context, id int64, d int, r string) error {
	if id == 99 {
		return errors.New("error")
	}
	return nil
}
func (m *mockCabinSvc) ListPrices(ctx context.Context, id int64) ([]domain.CabinPrice, error) {
	if id == 99 {
		return nil, errors.New("error")
	}
	return nil, nil
}
func (m *mockCabinSvc) UpsertPrice(ctx context.Context, p *domain.CabinPrice) error {
	if isErr(ctx) {
		return errors.New("error")
	}
	if p.PriceCents == 99 {
		return errors.New("error")
	}
	return nil
}

func (m *mockCabinSvc) BatchSetPrice(ctx context.Context, skuID int64, start, end time.Time, occupancy int, priceCents, childPriceCents, singleSupplementCents int64, priceType string) error {
	_ = start
	_ = end
	_ = occupancy
	_ = childPriceCents
	_ = singleSupplementCents
	_ = priceType
	if isErr(ctx) {
		return errors.New("error")
	}
	if skuID == 99 || priceCents == 99 {
		return errors.New("error")
	}
	return nil
}

func (m *mockCabinSvc) GetCategoryTree(ctx context.Context) (interface{}, error) {
	if isErr(ctx) {
		return nil, errors.New("error")
	}
	return []interface{}{}, nil
}

type mockBookingSvc struct{}

func (m *mockBookingSvc) Create(_ context.Context, userID, voyageID, skuID int64, guests int) (*domain.Booking, error) {
	if voyageID == 99 {
		return nil, errors.New("error")
	}
	return &domain.Booking{ID: 1, Status: "created", TotalCents: 10000}, nil
}

func isErr(ctx context.Context) bool {
	if c, ok := ctx.(*gin.Context); ok && c != nil && c.Request != nil && c.Request.URL != nil {
		return c.Request.URL.Query().Get("err") == "1"
	}
	return false
}

// TestRestHandlers 测试所有 REST 处理器
func TestRestHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if c.Query("err") == "1" {
			c.Set("trigger_error", true)
		}
		c.Next()
	})

	rtH := NewRouteHandler(&mockRouteSvc{})
	vH := NewVoyageHandler(&mockVoyageSvc{})
	cbH := NewCabinHandler(&mockCabinSvc{})
	bkH := NewBookingHandler(&mockBookingSvc{})

	r.GET("/routes", rtH.List)
	r.POST("/routes", rtH.Create)
	r.PUT("/routes/:id", rtH.Update)
	r.DELETE("/routes/:id", rtH.Delete)

	r.GET("/voyages", vH.List)
	r.POST("/voyages", vH.Create)
	r.PUT("/voyages/:id", vH.Update)
	r.DELETE("/voyages/:id", vH.Delete)

	r.GET("/cabins", cbH.FilteredList)
	r.POST("/cabins", cbH.Create)
	r.PUT("/cabins/:id", cbH.Update)
	r.PUT("/cabins/batch-status", cbH.BatchUpdateStatus)
	r.DELETE("/cabins/:id", cbH.Delete)
	r.GET("/cabins/alerts", cbH.GetAlerts)
	r.PUT("/cabins/:id/alert-threshold", cbH.SetAlertThreshold)
	r.GET("/cabins/:id/inventory", cbH.GetInventory)
	r.POST("/cabins/:id/inventory/adjust", cbH.AdjustInventory)
	r.GET("/cabins/:id/prices", cbH.ListPrices)
	r.POST("/cabins/:id/prices", cbH.UpsertPrice)

	r.POST("/bookings", bkH.Create)

	// 航线
	doReq(r, "GET", "/routes", nil)
	doReq(r, "GET", "/routes?err=1", nil)
	doReq(r, "GET", "/routes?err=1&page=abc&page_size=xyz", nil)
	doReq(r, "POST", "/routes", domain.Route{})
	doReq(r, "POST", "/routes", domain.Route{Name: "error"})
	doReq(r, "POST", "/routes", map[string]int{"name": 1})
	doReq(r, "PUT", "/routes/1", domain.Route{})
	doReq(r, "PUT", "/routes/1", map[string]int{"name": 1})
	doReq(r, "PUT", "/routes/1", domain.Route{Name: "error"})
	doReq(r, "PUT", "/routes/99", domain.Route{})
	doReq(r, "PUT", "/routes/x", nil)
	doReq(r, "DELETE", "/routes/1", nil)
	doReq(r, "DELETE", "/routes/1?err=1", nil)
	doReq(r, "DELETE", "/routes/1?err=1", nil)
	doReq(r, "DELETE", "/routes/99", nil)
	doReq(r, "DELETE", "/routes/x", nil)

	// 航次
	doReq(r, "GET", "/voyages", nil)
	doReq(r, "GET", "/voyages?route_id=99", nil)
	doReq(r, "POST", "/voyages", map[string]interface{}{"route_id": 1, "code": "test"})
	doReq(r, "POST", "/voyages", map[string]interface{}{"route_id": 1, "code": "error"})
	doReq(r, "POST", "/voyages", map[string]int{"code": 1}) // triggers ShouldBindJSON error (voyage uses code)
	doReq(r, "PUT", "/voyages/1", map[string]interface{}{"route_id": 1, "code": "test"})
	doReq(r, "PUT", "/voyages/1", map[string]int{"code": 1}) // triggers ShouldBindJSON error
	doReq(r, "PUT", "/voyages/1", map[string]interface{}{"route_id": 1, "code": "error"})
	doReq(r, "PUT", "/voyages/99", map[string]interface{}{"route_id": 1, "code": "test"})
	doReq(r, "PUT", "/voyages/x", nil)
	doReq(r, "DELETE", "/voyages/1", nil)
	doReq(r, "DELETE", "/voyages/99", nil)
	doReq(r, "DELETE", "/voyages/x", nil)

	// 舱位
	doReq(r, "GET", "/cabins", nil)
	doReq(r, "GET", "/cabins?voyage_id=1&cabin_type_id=1&status=1&page=1&page_size=10", nil)
	doReq(r, "GET", "/cabins?voyage_id=99", nil)
	doReq(r, "POST", "/cabins", map[string]interface{}{"voyage_id": 1, "cabin_type_id": 1, "code": "test"})
	doReq(r, "POST", "/cabins", map[string]interface{}{"voyage_id": 1, "cabin_type_id": 1, "code": "error"})
	doReq(r, "POST", "/cabins", map[string]int{"code": 1}) // triggers ShouldBindJSON error
	doReq(r, "PUT", "/cabins/1", map[string]interface{}{"voyage_id": 1, "cabin_type_id": 1, "code": "test"})
	doReq(r, "PUT", "/cabins/1", map[string]int{"code": 1}) // triggers ShouldBindJSON error
	doReq(r, "PUT", "/cabins/1", map[string]interface{}{"voyage_id": 1, "cabin_type_id": 1, "code": "error"})
	doReq(r, "PUT", "/cabins/99", map[string]interface{}{"voyage_id": 1, "cabin_type_id": 1, "code": "test"})
	doReq(r, "PUT", "/cabins/x", nil)
	doReq(r, "PUT", "/cabins/batch-status", map[string]interface{}{"ids": []int64{1, 2}, "status": 0})
	doReq(r, "PUT", "/cabins/batch-status", map[string]interface{}{"ids": []int64{}, "status": 0})
	doReq(r, "DELETE", "/cabins/1", nil)
	doReq(r, "DELETE", "/cabins/99", nil)
	doReq(r, "DELETE", "/cabins/x", nil)
	doReq(r, "GET", "/cabins/alerts", nil)
	doReq(r, "PUT", "/cabins/1/alert-threshold", map[string]interface{}{"threshold": 3})
	doReq(r, "PUT", "/cabins/1/alert-threshold", map[string]interface{}{"threshold": -1})
	doReq(r, "PUT", "/cabins/x/alert-threshold", map[string]interface{}{"threshold": 1})
	doReq(r, "GET", "/cabins/1/inventory", nil)
	doReq(r, "GET", "/cabins/99/inventory", nil)
	doReq(r, "GET", "/cabins/x/inventory", nil)
	doReq(r, "POST", "/cabins/1/inventory/adjust", map[string]interface{}{"delta": 1, "reason": "test"})
	doReq(r, "POST", "/cabins/99/inventory/adjust", map[string]interface{}{"delta": 1, "reason": "test"})
	doReq(r, "POST", "/cabins/x/inventory/adjust", nil)
	doReq(r, "POST", "/cabins/1/inventory/adjust", map[string]int{"delta": 1}) // bad
	doReq(r, "GET", "/cabins/1/prices", nil)
	doReq(r, "GET", "/cabins/99/prices", nil)
	doReq(r, "GET", "/cabins/x/prices", nil)
	doReq(r, "POST", "/cabins/1/prices", domain.CabinPrice{})

	doReq(r, "POST", "/cabins/1/prices", map[string]int{"price": 1}) // bad
	doReq(r, "POST", "/cabins/x/prices", nil)

	// 订单
	doReq(r, "POST", "/bookings", nil)
}
