package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/cruisebooking/backend/internal/handler"
	"github.com/gin-gonic/gin"
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
		Route:            &handler.RouteHandler{},
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

	// Test if a route is correctly registered
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/admin/auth/login", nil)
	router.ServeHTTP(w, req)

	// It should return some response, not 404
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}
