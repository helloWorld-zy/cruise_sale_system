package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	mw := JWT(JWTConfig{Secret: secret})

	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   "123",
		"roles": []interface{}{"admin", "user"},
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	validStr, _ := validToken.SignedString([]byte(secret))

	missingSubToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	missingSubStr, _ := missingSubToken.SignedString([]byte(secret))

	invalidAlgToken := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "123",
	})
	invalidAlgStr, _ := invalidAlgToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	tests := []struct {
		name   string
		header string
		status int
	}{
		{"No Header", "", http.StatusUnauthorized},
		{"Bad Prefix", "Token abc", http.StatusUnauthorized},
		{"Invalid Token", "Bearer invalid.token.str", http.StatusUnauthorized},
		{"Invalid Alg", "Bearer " + invalidAlgStr, http.StatusUnauthorized},
		{"Missing Sub", "Bearer " + missingSubStr, http.StatusUnauthorized},
		{"Valid Token", "Bearer " + validStr, http.StatusOK},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.Use(mw)
		r.GET("/", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})
		c.Request, _ = http.NewRequest("GET", "/", nil)
		if tt.header != "" {
			c.Request.Header.Set("Authorization", tt.header)
		}
		r.HandleContext(c)
		if w.Code != tt.status {
			t.Errorf("%s: expected %d, got %d", tt.name, tt.status, w.Code)
		}
	}
}
