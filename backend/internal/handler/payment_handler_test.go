package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakePaymentCallbackService struct{}

func (f fakePaymentCallbackService) HandleCallback(payload []byte) error { return nil }

func TestPaymentCallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewPaymentHandler(fakePaymentCallbackService{})
	r.POST("/api/pay/callback", h.Callback)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("POST", "/api/pay/callback", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
