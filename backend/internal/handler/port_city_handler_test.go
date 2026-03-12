package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type portCityServiceStub struct{}

func (s *portCityServiceStub) Search(_ context.Context, keyword string) ([]service.PortCityOption, error) {
	if keyword == "仁川" {
		return []service.PortCityOption{{Label: "仁川（韩国）", CityName: "仁川", CountryName: "韩国"}}, nil
	}
	return []service.PortCityOption{{Label: "海上巡游", IsSpecial: true}}, nil
}

func TestPortCityHandlerSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewPortCityHandler(&portCityServiceStub{})
	r.GET("/port-cities", h.Search)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/port-cities?keyword=仁川", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}
