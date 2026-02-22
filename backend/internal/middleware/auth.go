package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ContextKeyStaffID is the gin context key for the authenticated staff ID.
const ContextKeyStaffID = "staffID"

// ContextKeyRoles is the gin context key for the authenticated staff roles.
const ContextKeyRoles = "roles"

// JWTConfig holds JWT middleware configuration.
type JWTConfig struct {
	Secret string
}

// JWT validates the Bearer token, enforces HS256, and injects staffID + roles into context.
func JWT(cfg JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// CR-01 FIX: enforce HMAC signing method to prevent algorithm confusion attacks.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.Secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Extract staffID from "sub" claim.
		sub, err := claims.GetSubject()
		if err != nil || sub == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(ContextKeyStaffID, sub)

		// Extract roles from claims.
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
