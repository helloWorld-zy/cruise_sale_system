package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/handler"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type routerVoyageSvcStub struct{}

type routerContentTemplateSvcStub struct{}

type routerPortCitySvcStub struct{}

func (s *routerVoyageSvcStub) List(context.Context) ([]domain.Voyage, error) {
	return []domain.Voyage{}, nil
}
func (s *routerVoyageSvcStub) ListPublic(context.Context, int64, string, int, int) ([]domain.Voyage, int64, error) {
	return []domain.Voyage{}, 0, nil
}
func (s *routerVoyageSvcStub) Create(context.Context, *domain.Voyage) error { return nil }
func (s *routerVoyageSvcStub) Update(context.Context, *domain.Voyage) error { return nil }
func (s *routerVoyageSvcStub) GetByID(context.Context, int64) (*domain.Voyage, error) {
	return &domain.Voyage{ID: 1}, nil
}
func (s *routerVoyageSvcStub) Delete(context.Context, int64) error { return nil }

func (s *routerContentTemplateSvcStub) List(context.Context, domain.ContentTemplateKind) ([]domain.ContentTemplate, error) {
	return []domain.ContentTemplate{}, nil
}
func (s *routerContentTemplateSvcStub) GetByID(context.Context, int64) (*domain.ContentTemplate, error) {
	return &domain.ContentTemplate{ID: 1, Name: "默认费用说明", Kind: domain.ContentTemplateKindFeeNote, ContentJSON: "{}"}, nil
}
func (s *routerContentTemplateSvcStub) Create(context.Context, *domain.ContentTemplate) error {
	return nil
}
func (s *routerContentTemplateSvcStub) Update(context.Context, *domain.ContentTemplate) error {
	return nil
}
func (s *routerContentTemplateSvcStub) Delete(context.Context, int64) error { return nil }

func (s *routerPortCitySvcStub) Search(context.Context, string) ([]service.PortCityOption, error) {
	return []service.PortCityOption{}, nil
}

func TestSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Since we're just testing route registration, we can use empty handlers
	deps := Dependencies{
		Auth:             &handler.AuthHandler{},
		Company:          &handler.CompanyHandler{},
		Cruise:           &handler.CruiseHandler{},
		CabinType:        &handler.CabinTypeHandler{},
		FacilityCategory: &handler.FacilityCategoryHandler{},
		Facility:         &handler.FacilityHandler{},
		Image:            &handler.ImageHandler{},
		Voyage:           handler.NewVoyageHandler(&routerVoyageSvcStub{}),
		Cabin:            &handler.CabinHandler{},
		Booking:          &handler.BookingHandler{},
		User:             &handler.UserHandler{},
		Upload:           &handler.UploadHandler{},
		Payment:          &handler.PaymentHandler{},
		Refund:           &handler.RefundHandler{},
		Analytics:        &handler.AnalyticsHandler{},
		PortCity:         handler.NewPortCityHandler(&routerPortCitySvcStub{}),
		ContentTemplate:  handler.NewContentTemplateHandler(&routerContentTemplateSvcStub{}),
		JWTSecret:        "test-secret",
		Enforcer:         &casbin.Enforcer{},
	}

	router := Setup(deps)
	assert.NotNil(t, router)

	// Test if routes are correctly registered
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/admin/auth/login", nil)
	router.ServeHTTP(w, req)

	// It should return some response, not 404
	assert.NotEqual(t, http.StatusNotFound, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/api/v1/admin/cruises/batch-status", nil)
	router.ServeHTTP(w2, req2)
	assert.NotEqual(t, http.StatusNotFound, w2.Code)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/admin/images?entity_type=cruise&entity_id=1", nil)
	router.ServeHTTP(w3, req3)
	assert.NotEqual(t, http.StatusNotFound, w3.Code)

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("GET", "/api/v1/voyages?cruise_id=11&page=1&page_size=20", nil)
	router.ServeHTTP(w4, req4)
	assert.NotEqual(t, http.StatusNotFound, w4.Code)

	w5 := httptest.NewRecorder()
	req5, _ := http.NewRequest("GET", "/api/v1/voyages/101", nil)
	router.ServeHTTP(w5, req5)
	assert.NotEqual(t, http.StatusNotFound, w5.Code)

	w6 := httptest.NewRecorder()
	req6, _ := http.NewRequest("GET", "/api/v1/admin/content-templates?kind=fee_note", nil)
	router.ServeHTTP(w6, req6)
	assert.NotEqual(t, http.StatusNotFound, w6.Code)

	w7 := httptest.NewRecorder()
	req7, _ := http.NewRequest("GET", "/api/v1/admin/port-cities?keyword=仁川", nil)
	router.ServeHTTP(w7, req7)
	assert.NotEqual(t, http.StatusNotFound, w7.Code)
}

func TestSetup_ProtectedBatchEndpointsRequireAuthAndRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	modelPath := t.TempDir() + "/rbac_model.conf"
	policyPath := t.TempDir() + "/rbac_policy.csv"
	_ = os.WriteFile(modelPath, []byte(`
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`), 0644)
	_ = os.WriteFile(policyPath, []byte(`
p, admin, /api/v1/admin/cruises/batch-status, PUT
p, admin, /api/v1/admin/cabins/batch-status, PUT
`), 0644)
	enforcer, _ := casbin.NewEnforcer(modelPath, policyPath)

	deps := Dependencies{JWTSecret: "test-secret", Enforcer: enforcer, Auth: &handler.AuthHandler{}, Company: &handler.CompanyHandler{}, Cruise: &handler.CruiseHandler{}, CabinType: &handler.CabinTypeHandler{}, FacilityCategory: &handler.FacilityCategoryHandler{}, Facility: &handler.FacilityHandler{}, Image: &handler.ImageHandler{}, Voyage: handler.NewVoyageHandler(&routerVoyageSvcStub{}), Cabin: &handler.CabinHandler{}, Booking: &handler.BookingHandler{}, User: &handler.UserHandler{}, Upload: &handler.UploadHandler{}, Payment: &handler.PaymentHandler{}, Refund: &handler.RefundHandler{}, Analytics: &handler.AnalyticsHandler{}, PortCity: handler.NewPortCityHandler(&routerPortCitySvcStub{}), ContentTemplate: handler.NewContentTemplateHandler(&routerContentTemplateSvcStub{})}
	r := Setup(deps)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/cruises/batch-status", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "123", "exp": time.Now().Add(time.Hour).Unix()})
	signed, _ := token.SignedString([]byte("test-secret"))
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/cabins/batch-status", nil)
	req2.Header.Set("Authorization", "Bearer "+signed)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusForbidden, w2.Code)
}
