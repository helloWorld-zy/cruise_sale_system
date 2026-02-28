package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeRefundSvc struct{ err error }

func (f fakeRefundSvc) Create(_ context.Context, _ int64, _ int64, _ string) error { return f.err }

// TestRefundHandler_CreateOK 测试创建退款成功
func TestRefundHandler_CreateOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewRefundHandler(fakeRefundSvc{})
	r := gin.New()
	r.POST("/refunds", h.Create)

	body := bytes.NewBufferString(`{"payment_id":1,"amount_cents":5000,"reason":"cancel"}`)
	req := httptest.NewRequest("POST", "/refunds", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestRefundHandler_MissingPaymentID 测试缺少支付ID时的处理
func TestRefundHandler_MissingPaymentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewRefundHandler(fakeRefundSvc{})
	r := gin.New()
	r.POST("/refunds", h.Create)

	req := httptest.NewRequest("POST", "/refunds", bytes.NewBufferString(`{"amount_cents":500,"reason":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRefundHandler_MissingAmountCents 测试缺少金额时的处理
func TestRefundHandler_MissingAmountCents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewRefundHandler(fakeRefundSvc{})
	r := gin.New()
	r.POST("/refunds", h.Create)

	req := httptest.NewRequest("POST", "/refunds", bytes.NewBufferString(`{"payment_id":1,"reason":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRefundHandler_MissingReason 测试缺少原因时的处理
func TestRefundHandler_MissingReason(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewRefundHandler(fakeRefundSvc{})
	r := gin.New()
	r.POST("/refunds", h.Create)

	req := httptest.NewRequest("POST", "/refunds", bytes.NewBufferString(`{"payment_id":1,"amount_cents":500}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRefundHandler_ServiceError_ExceedsAmount 测试超过退款金额时的服务错误
func TestRefundHandler_ServiceError_ExceedsAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewRefundHandler(fakeRefundSvc{err: errors.New("refund amount exceeds remaining refundable balance 0")})
	r := gin.New()
	r.POST("/refunds", h.Create)

	body := bytes.NewBufferString(`{"payment_id":1,"amount_cents":99999,"reason":"too much"}`)
	req := httptest.NewRequest("POST", "/refunds", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
