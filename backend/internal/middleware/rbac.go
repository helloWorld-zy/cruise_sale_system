package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// RBAC 返回一个基于 Casbin 的角色访问控制中间件。
// 该中间件检查调用者的角色是否被 Casbin 策略允许访问当前请求的路径和方法。
// 角色信息从 gin 上下文中读取（由 JWT 中间件设置，键名为 "roles"）。
func RBAC(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// CR-02 修复：从 JWT 中间件设置的上下文中读取角色列表
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

		// 获取当前请求的资源路径和操作方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 检查是否有任一角色具有访问权限
		for _, role := range roles {
			allowed, err := enforcer.Enforce(role, obj, act)
			if err == nil && allowed {
				c.Next()
				return
			}
		}

		// 所有角色均无权限，拒绝访问
		c.AbortWithStatus(http.StatusForbidden)
	}
}
