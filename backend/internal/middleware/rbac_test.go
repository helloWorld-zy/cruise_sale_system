package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func TestRBAC(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Write temporary rbac config files
	modelPath := t.TempDir() + "/rbac_model.conf"
	policyPath := t.TempDir() + "/rbac_policy.csv"

	os.WriteFile(modelPath, []byte(`
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`), 0644)

	os.WriteFile(policyPath, []byte(`
p, admin, /api/secret, GET
`), 0644)

	enforcer, _ := casbin.NewEnforcer(modelPath, policyPath)
	mw := RBAC(enforcer)

	tests := []struct {
		name   string
		roles  interface{}
		path   string
		status int
	}{
		{"No Roles in Context", nil, "/api/secret", http.StatusForbidden},
		{"Empty Roles List", []string{}, "/api/secret", http.StatusForbidden},
		{"Invalid Role Type", "admin", "/api/secret", http.StatusForbidden},
		{"Role Denied", []string{"user"}, "/api/secret", http.StatusForbidden},
		{"Role Allowed", []string{"admin"}, "/api/secret", http.StatusOK},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(func(ctx *gin.Context) {
			if tt.roles != nil {
				ctx.Set(ContextKeyRoles, tt.roles)
			}
			ctx.Next()
		})
		r.Use(mw)

		r.GET("/api/secret", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})
		c.Request, _ = http.NewRequest("GET", "/api/secret", nil)
		r.HandleContext(c)

		if w.Code != tt.status {
			t.Errorf("%s: expected %d, got %d", tt.name, tt.status, w.Code)
		}
	}
}
