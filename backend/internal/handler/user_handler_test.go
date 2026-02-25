package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type userHandlerTestAuthSvc struct{ ok bool }

func (s userHandlerTestAuthSvc) VerifySMS(phone, code string) bool { return s.ok }
func (s userHandlerTestAuthSvc) SendSMS(_, _ string) error         { return nil }

type userHandlerSendErrAuthSvc struct{ err error }

func (s userHandlerSendErrAuthSvc) VerifySMS(phone, code string) bool { return true }
func (s userHandlerSendErrAuthSvc) SendSMS(_, _ string) error         { return s.err }

func TestUserHandlerLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewUserHandler(userHandlerTestAuthSvc{ok: true}, "secret")
	r.POST("/api/users/login", h.Login)

	w := httptest.NewRecorder()
	body := []byte(`{"phone":"13800000000","code":"1234"}`)
	req := httptest.NewRequest("POST", "/api/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestUserHandlerProfileRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewUserHandler(userHandlerTestAuthSvc{ok: true}, "secret")
	r.GET("/api/users/profile", h.Profile)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/users/profile", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestUserHandlerProfileSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, "13800000000") // M-01: C端使用 ContextKeyUserID
		c.Next()
	})

	h := NewUserHandler(userHandlerTestAuthSvc{ok: true}, "secret")
	r.GET("/api/users/profile", h.Profile)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/users/profile", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUserHandlerSendCodeSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewUserHandler(userHandlerTestAuthSvc{ok: true}, "secret")
	r.POST("/api/users/sms-code", h.SendCode)

	w := httptest.NewRecorder()
	body := []byte(`{"phone":"13800000000"}`)
	req := httptest.NewRequest("POST", "/api/users/sms-code", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestUserHandlerSendCodeTooFrequent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewUserHandler(userHandlerSendErrAuthSvc{err: service.ErrSMSTooFrequent}, "secret")
	r.POST("/api/users/sms-code", h.SendCode)

	w := httptest.NewRecorder()
	body := []byte(`{"phone":"13800000000"}`)
	req := httptest.NewRequest("POST", "/api/users/sms-code", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d, body=%s", w.Code, w.Body.String())
	}
}
