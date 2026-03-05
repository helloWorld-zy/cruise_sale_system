package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/cruisebooking/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

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
		Voyage:           &handler.VoyageHandler{},
		Cabin:            &handler.CabinHandler{},
		Booking:          &handler.BookingHandler{},
		User:             &handler.UserHandler{},
		Upload:           &handler.UploadHandler{},
		Payment:          &handler.PaymentHandler{},
		Refund:           &handler.RefundHandler{},
		Analytics:        &handler.AnalyticsHandler{},
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

	deps := Dependencies{JWTSecret: "test-secret", Enforcer: enforcer, Auth: &handler.AuthHandler{}, Company: &handler.CompanyHandler{}, Cruise: &handler.CruiseHandler{}, CabinType: &handler.CabinTypeHandler{}, FacilityCategory: &handler.FacilityCategoryHandler{}, Facility: &handler.FacilityHandler{}, Image: &handler.ImageHandler{}, Voyage: &handler.VoyageHandler{}, Cabin: &handler.CabinHandler{}, Booking: &handler.BookingHandler{}, User: &handler.UserHandler{}, Upload: &handler.UploadHandler{}, Payment: &handler.PaymentHandler{}, Refund: &handler.RefundHandler{}, Analytics: &handler.AnalyticsHandler{}}
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
