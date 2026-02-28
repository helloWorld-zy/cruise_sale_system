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

type fakePayCallbackSvc struct{ err error }

func (f fakePayCallbackSvc) HandleCallback(_ context.Context, _ string, _ []byte, _ string) error {
	return f.err
}

// TestPaymentCallback_OK 测试支付回调成功
func TestPaymentCallback_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPaymentHandler(fakePayCallbackSvc{})
	r := gin.New()
	r.POST("/callback", h.Callback)

	body := []byte(`{"trade_no":"TX001"}`)
	req := httptest.NewRequest("POST", "/callback?provider=wechat", bytes.NewReader(body))
	req.Header.Set("Wechatpay-Signature", "valid-sig")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SUCCESS")
}

// TestPaymentCallback_ServiceError 测试支付回调服务错误
func TestPaymentCallback_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPaymentHandler(fakePayCallbackSvc{err: errors.New("invalid signature")})
	r := gin.New()
	r.POST("/callback", h.Callback)

	body := []byte(`{"trade_no":"TX001"}`)
	req := httptest.NewRequest("POST", "/callback?provider=wechat", bytes.NewReader(body))
	req.Header.Set("Wechatpay-Signature", "badsig")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 支付平台总是接收 HTTP 200；失败信息在响应体中。
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "FAIL")
}

// TestPaymentCallback_EmptyBody 测试支付回调空请求体
func TestPaymentCallback_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPaymentHandler(fakePayCallbackSvc{})
	r := gin.New()
	r.POST("/callback", h.Callback)

	req := httptest.NewRequest("POST", "/callback", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "FAIL")
}

// TestPaymentCallback_NoSignature 测试支付回调无签名
func TestPaymentCallback_NoSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPaymentHandler(fakePayCallbackSvc{})
	r := gin.New()
	r.POST("/callback", h.Callback)

	body := []byte(`{"trade_no":"TX001"}`)
	req := httptest.NewRequest("POST", "/callback", bytes.NewReader(body))
	// 没有 Wechatpay-Signature 请求头且没有 sign 参数。
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "FAIL")
}

// TestPaymentCallback_DefaultProvider_Wechat 测试支付回调默认服务商微信
func TestPaymentCallback_DefaultProvider_Wechat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewPaymentHandler(fakePayCallbackSvc{})
	r := gin.New()
	r.POST("/callback", h.Callback)

	body := []byte(`{"trade_no":"TX001"}`)
	req := httptest.NewRequest("POST", "/callback", bytes.NewReader(body))
	req.Header.Set("Wechatpay-Signature", "any-sig")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SUCCESS")
}
