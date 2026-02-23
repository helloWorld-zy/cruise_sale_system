package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ContextKeyStaffID 是 gin 上下文中存储已认证员工 ID 的键名。
const ContextKeyStaffID = "staffID"

// ContextKeyRoles 是 gin 上下文中存储已认证员工角色列表的键名。
const ContextKeyRoles = "roles"

// JWTConfig 包含 JWT 中间件的配置参数。
type JWTConfig struct {
	Secret string // JWT 签名密钥
}

// JWT 返回一个 Gin 中间件函数，用于验证 Bearer 令牌。
// 该中间件强制使用 HS256 签名算法，并将员工 ID 和角色信息注入到请求上下文中。
func JWT(cfg JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization 字段
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 提取令牌字符串
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		// 解析并验证令牌
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// CR-01 修复：强制验证 HMAC 签名方法，防止算法混淆攻击
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.Secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 提取 claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 从 "sub" 声明中提取员工 ID
		sub, err := claims.GetSubject()
		if err != nil || sub == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(ContextKeyStaffID, sub)

		// 从 claims 中提取角色列表
		if rolesRaw, exists := claims["roles"]; exists {
			switch v := rolesRaw.(type) {
			case []interface{}:
				roles := make([]string, 0, len(v))
				for _, r := range v {
					if s, ok := r.(string); ok {
						roles = append(roles, s)
					}
				}
				c.Set(ContextKeyRoles, roles)
			}
		}

		c.Next()
	}
}
