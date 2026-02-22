package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testSecret = "test-secret-key"

func makeToken(staffID string, roles []string, secret string, exp time.Duration) string {
	claims := jwt.MapClaims{
		"sub":   staffID,
		"roles": roles,
		"exp":   time.Now().Add(exp).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := tok.SignedString([]byte(secret))
	return s
}

func requestWithAuth(token string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

// --- JWT middleware tests ---

func TestJWT_MissingHeader_Returns401(t *testing.T) {
	r := gin.New()
	r.GET("/", JWT(JWTConfig{Secret: testSecret}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestJWT_InvalidToken_Returns401(t *testing.T) {
	r := gin.New()
	r.GET("/", JWT(JWTConfig{Secret: testSecret}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, requestWithAuth("not.a.token"))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestJWT_WrongSecret_Returns401(t *testing.T) {
	tok := makeToken("1", []string{"admin"}, "other-secret", time.Hour)
	r := gin.New()
	r.GET("/", JWT(JWTConfig{Secret: testSecret}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, requestWithAuth(tok))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestJWT_ExpiredToken_Returns401(t *testing.T) {
	tok := makeToken("1", []string{"admin"}, testSecret, -time.Hour)
	r := gin.New()
	r.GET("/", JWT(JWTConfig{Secret: testSecret}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, requestWithAuth(tok))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestJWT_AlgorithmNone_Returns401(t *testing.T) {
	// Craft a token signed with "none" method — should be rejected.
	header := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0"      // {"alg":"none","typ":"JWT"}
	payload := "eyJzdWIiOiIxIiwicm9sZXMiOlsiYWRtaW4iXX0" // {"sub":"1","roles":["admin"]}
	noneToken := header + "." + payload + "."

	r := gin.New()
	r.GET("/", JWT(JWTConfig{Secret: testSecret}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, requestWithAuth(noneToken))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for alg=none, got %d", w.Code)
	}
}

func TestJWT_ValidToken_InjectsContextAndReturns200(t *testing.T) {
	tok := makeToken("42", []string{"admin", "editor"}, testSecret, time.Hour)
	r := gin.New()
	r.GET("/", JWT(JWTConfig{Secret: testSecret}), func(c *gin.Context) {
		staffID, _ := c.Get(ContextKeyStaffID)
		roles, _ := c.Get(ContextKeyRoles)
		if staffID == nil || roles == nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, requestWithAuth(tok))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
