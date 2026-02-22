package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// RBAC checks that at least one of the caller's roles is permitted by Casbin policy.
// It reads roles from gin context set by the JWT middleware (key: "roles").
func RBAC(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// CR-02 FIX: read roles from context set by JWT middleware.
		rolesVal, exists := c.Get(ContextKeyRoles)
		if !exists {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		roles, ok := rolesVal.([]string)
		if !ok || len(roles) == 0 {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		obj := c.Request.URL.Path
		act := c.Request.Method

		for _, role := range roles {
			allowed, err := enforcer.Enforce(role, obj, act)
			if err == nil && allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatus(http.StatusForbidden)
	}
}
